package main

import (
	"context"
	"fmt"
	"io/ioutil" // Handling Files
	"log"
	"os"
	"strconv" // Converting Number to String

	//"encoding/json"
	"net"
	//"lab1/game/etapas"
	pb "lab1/game/helloworld"

	"google.golang.org/grpc"
)

type DataNode struct {
	pb.UnimplementedGameServer
}

//Int32 to String
func ItoS(n int32) string {
	buf := [11]byte{}
	pos := len(buf)
	i := int64(n)
	signed := i < 0
	if signed {
		i = -i
	}
	for {
		pos--
		buf[pos], i = '0'+byte(i%10), i/10
		if i == 0 {
			if signed {
				pos--
				buf[pos] = '-'
			}
			return string(buf[pos:])
		}
	}
}

func createFile(path string) {

	var _, err = os.Stat(path)

	if os.IsNotExist(err) {

		var file, err = os.Create(path)

		if err != nil {
			log.Fatalf("Error al crear archivo: %v", err)
			return
		}

		defer file.Close()

	}

}

const (
	port = ":50053"
)

func (s *DataNode) SendJugadas(ctx context.Context, in *pb.StagePlays) (*pb.StoredReply, error) {

	id_jugador := in.GetId()
	etapa := in.GetEtapa()
	jugadas := in.GetJugadasList()

	path := "files/jugador_" + ItoS(id_jugador) + "__ronda_" + ItoS(etapa) + ".txt"

	data := ""

	for _, jugada := range jugadas {
		fmt.Println(jugada)
		data += ItoS(jugada) + "\n"
	}

	dataByte := []byte(data)

	// Values Ronda
	err := ioutil.WriteFile(path, dataByte, 0777)

	var message string

	if err != nil {
		log.Fatalf("Error al registrar jugada: %v", err)
		message = "Error al registrar jugada"
	} else {
		message = "Jugador registrado con exito"
	}

	return &pb.StoredReply{Message: message}, nil

}

func main() {
	// Puerto para escucha de los clientes
	lis, err := net.Listen("tcp", port)

	if err != nil {
		log.Fatalf("Error al escuchar: %v", err)
	}

	// Crea una isntancia de server gRPC
	s := grpc.NewServer()

	// Registrar nuestra implementación de server con gRPC server
	pb.RegisterGameServer(s, &DataNode{})

	log.Printf("Server escuchando en %v", lis.Addr())
	// Llamar Server() con los detalles de puerto para realizar un bloquero
	// de espera hasta que el proceso sea killed o Stop() es llamado
	if err := s.Serve(lis); err != nil {
		log.Fatalf("failed to serve: %v", err)
	}

	pdata := []byte("All my data i want to write to a file")

	err = ioutil.WriteFile("player.txt", pdata, 0777)

	if err != nil {
		fmt.Println(err)
	}

	// Removing Files Players
	for r := 0; r < 10; r++ {
		for i := 0; i < 10; i++ {
			err = os.Remove("jugador_" + strconv.Itoa(i) + "__ronda_" + strconv.Itoa(r) + ".txt")

			if err != nil {
				log.Fatal(err)
			}

		}
	}

}
