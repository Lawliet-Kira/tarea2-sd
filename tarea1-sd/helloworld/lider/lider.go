package main

import (
	"context"
	//"encoding/json"
	"fmt"
	"log"
	"math"
	"math/rand"
	"net"
	"time"

	//"lab1/game/etapas"
	pb "lab1/game/helloworld"

	"google.golang.org/grpc"
)

type Jugada struct {
	id     int32
	jugada int32
}

var jugadoresOnline []int32
var contador int32 = 0
var count int32 = 1
var etapa = 1
var jugadores_etapa []int32
var jugadores_ronda []int32
var jugadores_turno []int32
var jugadores_E2 []int32
var jugadores_E3 []Jugada
var begin bool = false
var shuffled bool = false
var maxplayers int32 = 2
var players_alive int32 = 0
var players_round int32 = 0
var optLider int32
var beginR bool = true
var contador_muertos int32 = 0
var equipos [2][]int32
var equipo_ganador int32 = 0
var equiposE3 [][]int32
var suma_T1 int32 = 0
var suma_T2 int32 = 0
var murio int32
var cant_equiposE3 int

var separados bool = false

// Variables Globales etapa 1

var roundE1 int32 = 0

/*
		* Enviar señal de inicio de etapa
		* Actualizar sobre el estado en el juego
	-> Pozo
		* Enviar jugadores eliminados
		* Pedir monto total acumulado
	-> NameNode
		* Enviar jugadas de los jugadores

- Recepción:
	->	Jugadores
		* Participar Juego
		* Enviar jugadas
	-> Pozo
		* Devolver monto total del pozo
	-> NameNode
		*Devolver registro del jugador
*/

// server is used to implement helloworld.GreeterServer.
type server struct {
	pb.UnimplementedGameServer
}

// Etapas
/*func Etapa1() {

	// Mensaje comienza etapa 1 a jugadores
	total_jug_vivos := 0
	total_jug_resp := 0

	for total_jug_resp < total_jug_vivos {


	}

	for rondas < 3 {

		total_jug_resp = 0

		// Recibir jugadas
		for total_jug_resp < total_jug_vivos {

		}

		// Respuesta del lider
		numLider := etapas.RandomNumberLider()

		for j := 0;  j < total_jug_vivos; j++ {

		}

		if numLider >    {


		}



	}

	return 1
}


func Etapa2() {

	// Dividir jugadores en equipo

	for round :=  ; round <


	numLider :=	etapa.RandomNumberStageTwo()


}*/
//Utility
func contains(s []int32, elem int32) bool {
	for _, v := range s {
		if v == elem {
			return true
		}
	}
	return false
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

func (s *server) JoinGame(ctx context.Context, in *pb.JoinRequest) (*pb.JoinReply, error) {
	log.Printf("Received: %v", in.GetMessage())

	var id int32
	var message string

	if in.GetId() == 0 { // Jugador no seteado
		id = count
		count++
		players_alive = players_alive + 1
		jugadoresOnline = append(jugadoresOnline, id)
		message = "Te has inscrito al juego del calamar, Bienvenido Jugador" + ItoS(id)
	} else {
		id = in.GetId()
		message = "Esperando al resto de jugadores..."
	}

	log.Printf("Jugador añadido ")
	fmt.Printf("%v", jugadoresOnline)
	return &pb.JoinReply{Id: id, Message: message}, nil
}

// BeginGame implements game.GameServer
func (s *server) BeginGame(ctx context.Context, in *pb.BeginRequest) (*pb.BeginReply, error) {
	log.Printf("Recibido de "+ItoS(in.GetId())+": %v", in.GetMessage())
	begin = false
	if players_alive == maxplayers {
		begin = true
		log.Printf("Comienza el juego")
	}
	return &pb.BeginReply{Message: begin, Online: int32(players_alive)}, nil
}

//BeginStage
func (s *server) BeginStage(ctx context.Context, in *pb.BeginStageRequest) (*pb.BeginStageReply, error) {

	rand.Seed(time.Now().UnixNano())
	begin = false
	id_jugador := in.GetId()

	log.Printf("Received from "+ItoS(id_jugador)+": %v", in.GetStage())
	contiene := contains(jugadores_etapa, id_jugador)

	if !contiene {
		jugadores_etapa = append(jugadores_etapa, id_jugador)
	}
	fmt.Println(players_alive, jugadores_etapa, etapa)
	if int32(len(jugadores_etapa)) == players_alive {
		begin = true
	}
	if roundE1 == 4 {
		etapa = 2
	}
	if etapa == 2 && begin {
		if !shuffled {
			fmt.Println(jugadores_etapa)
			//shuffle
			for i := range jugadores_etapa {
				j := rand.Intn(i + 1)
				jugadores_etapa[i], jugadores_etapa[j] = jugadores_etapa[j], jugadores_etapa[i]
			}
			shuffled = true
			fmt.Println(jugadores_etapa)
			optLider = int32(rand.Intn(4) + 1)
			if len(jugadores_etapa)%2 == 1 {
				murio = jugadores_etapa[len(jugadores_etapa)-1]
				jugadores_etapa = jugadores_etapa[:len(jugadores_etapa)-1]
				players_alive--

			}

			equipos[0] = jugadores_etapa[:len(jugadores_etapa)/2]
			equipos[1] = jugadores_etapa[len(jugadores_etapa)/2:]

			fmt.Println(equipos)
		}
		if id_jugador == murio {
			return &pb.BeginStageReply{Stage: begin, Dead: true}, nil
		}
	}
	if etapa == 3 && begin {
		if !separados {
			fmt.Println(jugadores_etapa)
			for i := range jugadores_etapa {
				j := rand.Intn(i + 1)
				jugadores_etapa[i], jugadores_etapa[j] = jugadores_etapa[j], jugadores_etapa[i]
			}
			if len(jugadores_etapa)%2 == 1 {
				murio = jugadores_etapa[len(jugadores_etapa)-1]
				jugadores_etapa = jugadores_etapa[:len(jugadores_etapa)-1]
				players_alive--
			}
			cant_equiposE3 = len(jugadores_etapa) / 2
			for i := 0; i < cant_equiposE3; i++ {
				equiposE3[i] = jugadores_etapa[2*i : 2*(i+1)]
			}
			optLider = int32(rand.Intn(10) + 1)
			separados = true
			fmt.Println(equiposE3)
		}
		if id_jugador == murio {
			return &pb.BeginStageReply{Stage: begin, Dead: true}, nil
		}
	}

	return &pb.BeginStageReply{Stage: begin}, nil
}

// EndGame implements game.GameServer
func (s *server) EndGame(ctx context.Context, in *pb.EndRequest) (*pb.EndReply, error) {
	log.Printf("Received: %v", in.GetMessage())
	return &pb.EndReply{Message: "Hello " + in.GetMessage()}, nil
}

func (s *server) SendJugadaE1(ctx context.Context, in *pb.JugadaE1) (*pb.PlayerStatusE1, error) {

	var dead bool

	var curr_sum int32 = 0

	if beginR {
		//Ver la seed
		rand.Seed(time.Now().UnixNano())
		optLider = int32(rand.Intn(4) + 6)
		beginR = false
		jugadores_etapa = nil
		jugadores_ronda = nil
		jugadores_turno = nil
		roundE1++
		fmt.Println("El numero elegido por el lider es: " + ItoS(optLider))

	}
	curr_sum = in.GetJugada() + in.GetSumaActual()

	if in.GetJugada() >= optLider {
		dead = true
	} else {
		dead = false
	}

	if roundE1 == 4 && curr_sum < 21 {
		dead = true

	} else {
		dead = false
	}

	id_jugador := in.GetId()
	contiene := contains(jugadores_ronda, id_jugador)

	if !contiene {
		jugadores_ronda = append(jugadores_ronda, id_jugador)
	}

	if int32(len(jugadores_ronda)) == players_alive {
		beginR = true
		players_alive = players_alive - contador_muertos
		contador_muertos = 0

		fmt.Println("Comienza la siguiente ronda")
	}

	if dead {
		contador_muertos++
		fmt.Println("Ha muerto el jugador " + ItoS(in.GetId()))
	}
	conn, err := grpc.Dial(gamenodeAddress, grpc.WithInsecure(), grpc.WithBlock())

	if err != nil {
		log.Fatalf("did not connect: %v", err)
	}

	defer conn.Close()

	// Client Stub to perform RPCs
	client := pb.NewGameClient(conn)

	//ctx, cancel := context.WithTimeout(context.Background(), time.Second)
	r, err := client.SendJugadas(ctx, &pb.StagePlays{Id: id_jugador, Etapa: int32(etapa), Jugada: in.GetJugada()})
	fmt.Println(r.GetMessage())

	return &pb.PlayerStatusE1{SumaTotal: curr_sum, Dead: dead, Ronda: roundE1}, nil

}

func (s *server) SendJugadaE2(ctx context.Context, in *pb.JugadaE2) (*pb.PingMsg, error) {

	id_jugador := in.GetId()
	jugada_jugador := in.GetJugada()
	fmt.Println("El lider Ha elegido: " + ItoS(optLider) + " Jugador " +
		ItoS(id_jugador) + " ha elegido " + ItoS((jugada_jugador)))
	if !contains(jugadores_E2, id_jugador) {
		jugadores_E2 = append(jugadores_E2, id_jugador)
		if contains(equipos[0], id_jugador) {
			suma_T1 = suma_T1 + jugada_jugador
		} else {
			suma_T2 = suma_T2 + jugada_jugador
		}
	}

	conn, err := grpc.Dial(gamenodeAddress, grpc.WithInsecure(), grpc.WithBlock())

	if err != nil {
		log.Fatalf("did not connect: %v", err)
	}

	defer conn.Close()

	// Client Stub to perform RPCs
	client := pb.NewGameClient(conn)

	//ctx, cancel := context.WithTimeout(context.Background(), time.Second)
	r, err := client.SendJugadas(ctx, &pb.StagePlays{Id: id_jugador, Etapa: int32(etapa), Jugada: in.GetJugada()})
	fmt.Println(r.GetMessage())
	return &pb.PingMsg{Id: id_jugador}, nil

}

/* func (s *server) SendJugadaE3(ctx context.Context, in *pb.JugadaE3) (*pb.PingMsg, error) {

	id_jugador := in.GetId()
	jugada_jugador := in.GetJugada()
	fmt.Println("El lider Ha elegido: " + ItoS(optLider) + " Jugador " +
	ItoS(id_jugador) +" ha elegido " + ItoS((jugada_jugador)))
	if !contains(jugadores_E3, id_jugador) {
		Jugada.id = id_jugador
		Jugada.jugada = jugada_jugador
		jugadores_E3 = append(jugadores_E3, Jugada)
	}
	return &pb.PingMsg{Id: id_jugador}, nil

} */

func (s *server) IsAlready(ctx context.Context, in *pb.PingMsg) (*pb.PlayerStatusE2, error) {

	tie := false
	ready := false
	id_jugador := in.GetId()
	if len(jugadores_E2) == int(players_alive) {

		fmt.Println("Lider jugó :", optLider)

		if !(suma_T1%2 == optLider%2) && !(suma_T2%2 == optLider%2) {
			fmt.Println("")
			tie = true
		}

		if !(suma_T1%2 == optLider%2) && !tie {
			if contains(equipos[0], id_jugador) {
				fmt.Println("Jugador " + ItoS(id_jugador) + " Eliminado")
				jugadores_etapa = jugadores_etapa[:len(jugadores_etapa)/2]
				players_alive--
				return &pb.PlayerStatusE2{Ready: true, Dead: true}, nil
			}
		}
		if !(suma_T2%2 == optLider%2) && !tie {
			if contains(equipos[1], id_jugador) {
				fmt.Println("Jugador " + ItoS(id_jugador) + " Eliminado")
				jugadores_etapa = jugadores_etapa[len(jugadores_etapa)/2:]
				players_alive--
				return &pb.PlayerStatusE2{Ready: true, Dead: true}, nil
			}
		}
		ready = true
		begin = true
		etapa = 3
	}
	return &pb.PlayerStatusE2{Ready: ready, Dead: false}, nil

}

func buscarJugada(id_jugador int32) int32 {
	var jugada int32 = 0
	for _, i := range jugadores_E3 {
		if i.id == id_jugador {
			jugada = i.jugada
		}
	}
	return jugada
}

func (s *server) IsAlready2(ctx context.Context, in *pb.PingMsg) (*pb.PlayerStatusE3, error) {
	var dead bool = false
	ready := false
	var id_pareja int32 = 0
	var jugada_pareja int32 = 0
	var jugada_jugador int32 = 0
	id_jugador := in.GetId()
	if len(jugadores_E3) == int(players_alive) {
		fmt.Println("Lider jugó :", optLider)
		for i := range equiposE3 {
			if contains(equiposE3[i], id_jugador) {
				if equiposE3[i][0] == id_jugador {
					id_pareja = equiposE3[i][1]
				} else {
					id_pareja = equiposE3[i][0]
				}
			}
		}

		jugada_jugador = buscarJugada(id_jugador)
		jugada_pareja = buscarJugada(id_pareja)

		valor_absolutoJugador := int32(math.Abs(float64(optLider - jugada_jugador)))
		valor_absolutoPareja := int32(math.Abs(float64(optLider - jugada_pareja)))

		if valor_absolutoJugador == valor_absolutoPareja {
			ready = true
			dead = false
		}
		if valor_absolutoJugador < valor_absolutoPareja {
			ready = true
			dead = true
		}
		if valor_absolutoJugador > valor_absolutoPareja {
			ready = true
			dead = false
		}

	}
	return &pb.PlayerStatusE3{Ready: ready, Dead: dead}, nil

}

//func (s *server) SendJugadaE2(ctx context.Context, in *pb.JugadaE2)  (*pb.PingMsg, error){}
func (s *server) BeginRound(ctx context.Context, in *pb.PingMsg) (*pb.BeginRoundReply, error) {
	begin = false
	id_jugador := in.GetId()
	log.Printf("Received from : " + ItoS(id_jugador))
	contiene := contains(jugadores_turno, id_jugador)

	if !contiene {
		jugadores_turno = append(jugadores_turno, id_jugador)
	}

	if int32(len(jugadores_turno)) == players_alive {
		begin = true

	}
	return &pb.BeginRoundReply{Round: begin}, nil
}

const (
	port            = ":50051"
	gamenodeAddress = "10.6.43.114:50052"
)

func main() {

	end := false

	// Puerto para escucha de los clientes
	lis, err := net.Listen("tcp", port)

	if err != nil {
		log.Fatalf("Error al escuchar: %v", err)
	}

	// Crea una isntancia de server gRPC
	s := grpc.NewServer()

	// Registrar nuestra implementación de server con gRPC server
	pb.RegisterGameServer(s, &server{})

	log.Printf("Server escuchando en %v", lis.Addr())

	// Llamar Server() con los detalles de puerto para realizar un bloquero
	// de espera hasta que el proceso sea killed o Stop() es llamado
	if err := s.Serve(lis); err != nil {
		log.Fatalf("failed to serve: %v", err)
	}

	for !end {

		// Menu de inicio
		fmt.Println("Bienvenido a Squid Game.")
		fmt.Println("Selecciona una opción:")
		fmt.Println(" 1. Iniciar Juego")
		fmt.Println(" 2. Salir del Juego")

		var uselect string

		fmt.Scanln(&uselect)

		if uselect == "1" { // Inicia Juego

			for !begin {
				fmt.Println("Esperando Jugadores... [%d]", players_alive)
			}

			fmt.Println("Comienzo de la Etapa 1...")
			time.Sleep(2 * time.Second)
			//etapaUno()
			fmt.Println("Comienzo de la Etapa 2...")
			time.Sleep(2 * time.Second)
			//etapaDos()
			fmt.Println("Comienzo de la Etapa 3...")
			time.Sleep(2 * time.Second)
			//etapaTres()

		} else if uselect == "2" { // Salir del Juego
			end = true
			s.Stop()
		} else {

			fmt.Println("Ingresaste una opción no válida, intentalo denuevo.")

		}

	}

	return

}
