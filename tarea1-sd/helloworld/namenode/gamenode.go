package main

import (
	"context"
	"fmt"
	//"encoding/json"
	"log"
	"math/rand"
	"net"
	"time"

	//"lab1/game/etapas"
	pb "lab1/game/helloworld"

	"google.golang.org/grpc"
)

type NameNode struct {
	pb.UnimplementedGameServer
}

type Location struct {
	
	etapa int32
	jugadas []int32
	direccion string 
}

type Player struct{
	id int32
	played []Location
	
}
const (
	port = ":50052"
)
func newLocation(etapa int32, direccion string) * Location{
	var jugadas []int32
	new := Location{etapa: etapa, direccion: direccion, jugadas: jugadas}
	return &new
}

func newPlayer(id int32) * Player{
	var locationList []Location
	locationList = append(locationList, *newLocation(1, ""))
	locationList = append(locationList, *newLocation(2, ""))
	locationList = append(locationList, *newLocation(3, ""))
	new := Player{id: id, played: locationList}
	return &new
}

func SearchPlayerPlays(PlayerList []Player, id int32) []Location{
	var locating []Location
	locating = nil 
	for _,player := range PlayerList{
		if player.id == id{
			locating = player.played
			return locating
		}
	}
	return locating

}

var PlayerPlays []Player

func (s *NameNode) SendJugadas(ctx context.Context, in *pb.StagePlays) (*pb.StoredReply, error) {
	rand.Seed(time.Now().UnixNano())
	var locating string
	switch machine := rand.Intn(4); machine{
		case 1:
			locating = ""
		case 2:
			locating = ""
		case 3:
			locating = ""	
	}
	locating = "localhost:50053"
	if SearchPlayerPlays(PlayerPlays, in.GetId()) == nil{
		newVisitor := newPlayer(in.GetId())
		PlayerPlays = append(PlayerPlays, *newVisitor)
		fmt.Print("New Visitor  ")
		fmt.Println(in.GetId())
	}

	arr_locations := SearchPlayerPlays(PlayerPlays, in.GetId())

	//for _, v := range SearchPlayerPlays(PlayerPlays, in.GetId()){
	for i := 0 ; i < len(arr_locations) ; i++ {
		fmt.Print("Player:")
		fmt.Println(in.GetId())
		fmt.Print("Iterando en etapa ")
		fmt.Println(arr_locations[i].etapa)
		fmt.Print("Recibido ")
		fmt.Println(in.GetEtapa())
		if in.GetEtapa() == arr_locations[i].etapa {
			arr_locations[i].jugadas = append(arr_locations[i].jugadas, in.GetJugada())
			arr_locations[i].direccion = locating
			if in.GetEtapa() == 1 {
				fmt.Print("Jugadas:")
				fmt.Println(arr_locations[i].jugadas)
				if len(arr_locations[i].jugadas) == 4{
					fmt.Println("Enviando etapa 1 a datanode")
					conn, err := grpc.Dial(locating, grpc.WithInsecure(), grpc.WithBlock())

					if err != nil {
						log.Fatalf("did not connect: %v", err)
					}

					defer conn.Close()

					// Client Stub to perform RPCs
					client := pb.NewGameClient(conn)
					ctx := context.Background()
					//ctx, cancel := context.WithTimeout(context.Background(), time.Second)
					r, err := client.SendJugadas(ctx, &pb.StagePlays{Id: in.GetId(), Etapa: in.GetEtapa(), JugadasList: arr_locations[i].jugadas})
					return &pb.StoredReply{Message: r.GetMessage()}, nil
				}

			}else{
				//enviar a Datanode
				conn, err := grpc.Dial(locating, grpc.WithInsecure(), grpc.WithBlock())

				if err != nil {
					log.Fatalf("did not connect: %v", err)
				}

				defer conn.Close()

				// Client Stub to perform RPCs
				client := pb.NewGameClient(conn)
				ctx := context.Background()
				//ctx, cancel := context.WithTimeout(context.Background(), time.Second)
				r, err := client.SendJugadas(ctx, &pb.StagePlays{Id: in.GetId(), Etapa: in.GetEtapa(), JugadasList: arr_locations[i].jugadas})
				return &pb.StoredReply{Message: r.GetMessage()}, nil
			}
			
		}

	}
	return &pb.StoredReply{Message: "No se pudo guardar"}, nil
}

func (s *NameNode) PullJugadas(ctx context.Context, in *pb.PingMsg) (*pb.StagePlays, error){
	for _, v:= range SearchPlayerPlays(in.GetId()){
		conn, err := grpc.Dial(v.direccion, grpc.WithInsecure(), grpc.WithBlock())

			if err != nil {
				log.Fatalf("did not connect: %v", err)
			}

			defer conn.Close()

			// Client Stub to perform RPCs
			client := pb.NewGameClient(conn)
			ctx := context.Background()
			r, err := client.PullJugadas(ctx, &pb.PingMsg{Id: in.GetId()})
		}
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
	pb.RegisterGameServer(s, &NameNode{})

	log.Printf("Server escuchando en %v", lis.Addr())

	// Llamar Server() con los detalles de puerto para realizar un bloquero
	// de espera hasta que el proceso sea killed o Stop() es llamado
	if err := s.Serve(lis); err != nil {
		log.Fatalf("failed to serve: %v", err)
	}
}