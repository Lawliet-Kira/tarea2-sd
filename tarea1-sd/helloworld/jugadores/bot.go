package main

import (
	"context"
	"fmt"
	pb "lab1/game/helloworld"
	"log"
	"math/rand"
	"time"

	"google.golang.org/grpc"
)

const (
	address = "10.6.43.113:50051"
)

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

func Bot() {

	// Crear un gRPC canal para comunicarse con el servidor
	// 	-> Esto se crea pasando server address y port number a grpc.Dial()
	conn, err := grpc.Dial(address, grpc.WithInsecure(), grpc.WithBlock())
	rand.Seed(time.Now().UnixNano())
	if err != nil {
		log.Fatalf("did not connect: %v", err)
	}

	defer conn.Close()

	// Client Stub to perform RPCs
	client := pb.NewGameClient(conn)
	message := "HOLA deseo unirme a the game"
	// Contact the server and print out its response.
	name := "Jugador 1"
	var id int32 = 0
	ctx := context.Background()
	//ctx, cancel := context.WithTimeout(context.Background(), time.Second)
	r, err := client.JoinGame(ctx, &pb.JoinRequest{Id: id, Message: message, Name: name})

	id = r.GetId()

	if err != nil {
		log.Fatalf("could not greet: %v", err)
	}
	log.Printf("Jugador: %s", r.GetMessage())

	comienzaeljuego := false

	for !comienzaeljuego {
		time.Sleep(5 * time.Second)
		message := "Ya comienza el juego??"
		r, err := client.BeginGame(ctx, &pb.BeginRequest{Id: id, Message: message})
		if err != nil {
			log.Fatalf("could not greet: %v", err)
		}
		comienzaeljuego = r.GetMessage()
		Jugadoresonline := r.GetOnline()
		fmt.Println(">> Jugadores online: " + itoS(Jugadoresonline))
	}

	// SOLICITUD A LA PRIMERA ETAPA
	stage := "Va a comenzar la etapa??"
	comienzalaetapa := false

	for !comienzalaetapa {

		time.Sleep(10 * time.Second)
		r, err := client.BeginStage(ctx, &pb.BeginStageRequest{Id: id, Stage: stage})
		if err != nil {
			log.Fatalf("could not greet: %v", err)
		}
		comienzalaetapa = r.GetStage()
	}

	fmt.Println("Comienza Squid Game")
	//COMIENZA LA PRIMERA ETAPA
	var suma_actual int32 = 0
	dead := false
	var newEtapa bool = true
	// Loop de rondas
	for newEtapa {

		var numero int
		numero = rand.Intn(11)

		jugada := int32(numero)

		r, err := client.SendJugadaE1(ctx, &pb.JugadaE1{Id: id, Jugada: jugada, SumaActual: suma_actual})
		if err != nil {
			log.Fatalf("could not greet: %v", err)
		}
		dead = r.GetDead()
		if dead {
			fmt.Println("Has muerto...")
			return
		}
		suma_actual = r.GetSumaTotal()
		ronda := r.GetRonda()
		fmt.Println("Suma total" + itoS(suma_actual) + " ronda: " + itoS(ronda))
		if ronda == 4 {
			newEtapa = false
		}
		var newRound bool = false
		for !newRound {
			fmt.Println("Aun no comienza la ronda")
			time.Sleep(5 * time.Second)
			response, err := client.BeginRound(ctx, &pb.PingMsg{Id: id})
			if err != nil {
				log.Fatalf("could not greet: %v", err)
			}
			newRound = response.GetRound()

		}

	}
	var action string

	comienzalaetapa = false
	for !comienzalaetapa {

		time.Sleep(5 * time.Second)
		r, err := client.BeginStage(ctx, &pb.BeginStageRequest{Id: id, Stage: stage})
		if err != nil {
			log.Fatalf("could not greet: %v", err)
		}
		comienzalaetapa = r.GetStage()
		if !comienzalaetapa {
			fmt.Println("Registrado para la siguiente etapa, esperando jugadores")
		}

	}

	fmt.Println("Comienza la Etapa 2.")
	newEtapa = false
	for !newEtapa {

		var numero int

		numero = rand.Intn(5)
		jugada := int32(numero)

		_, err := client.SendJugadaE2(ctx, &pb.JugadaE2{Id: id, Jugada: jugada})
		if err != nil {
			log.Fatalf("could not greet: %v", err)
		}

		var newRound bool = false
		for !newRound {
			fmt.Println("Aun no comienza la ronda")
			time.Sleep(5 * time.Second)
			response, err := client.IsAlready(ctx, &pb.PingMsg{Id: id})
			if err != nil {
				log.Fatalf("could not greet: %v", err)
			}
			newRound = response.GetReady()
			newEtapa = true
			dead = response.GetDead()
			if dead {
				return
			}
		}

	}

	comienzalaetapa = false
	for !comienzalaetapa {
		time.Sleep(5 * time.Second)
		r, err := client.BeginStage(ctx, &pb.BeginStageRequest{Id: id, Stage: stage})
		if err != nil {
			log.Fatalf("could not greet: %v", err)
		}
		comienzalaetapa = r.GetStage()
		if !comienzalaetapa {
			fmt.Println("Registrado para la siguiente etapa, esperando jugadores")
		}

	}

	newEtapa = false
	for !newEtapa {

		var numero int
		numero = rand.Intn(10)
		jugada := int32(numero)

		_, err := client.SendJugadaE3(ctx, &pb.JugadaE3{Id: id, Jugada: jugada})
		if err != nil {
			log.Fatalf("could not greet: %v", err)
		}

		var newRound bool = false
		for !newRound {
			fmt.Println("Aun no comienza la ronda")
			time.Sleep(5 * time.Second)
			response, err := client.IsAlready2(ctx, &pb.PingMsg{Id: id})
			if err != nil {
				log.Fatalf("could not greet: %v", err)
			}
			newRound = response.GetReady()
			newEtapa = true
			dead = response.GetDead()
			if dead {
				return
			}
		}
	}
	return

}

func main() {

	// Crear un gRPC canal para comunicarse con el servidor
	// 	-> Esto se crea pasando server address y port number a grpc.Dial()
	for i := 1; i < 16; i++ {
		go Bot()
	}

	return

}
