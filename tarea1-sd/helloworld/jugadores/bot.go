package main

import (
	"context"
	"log"
	"time"
	"fmt"
	pb "lab1/game/helloworld"
)

const (
	address     = "localhost:50051"
	defaultName = "world"
)

func Bot(ctx context.Context,client pb.GameClient) {

	message := "Hola soy el jugador x"
	name := "Jugador Bot"

	var id int32 = 0
	r, err := client.JoinGame(ctx, &pb.JoinRequest{Id: id, Message: message, Name: name})
	
	id = r.GetId()
	
	if err != nil {
		log.Fatalf("Error no pudo unirse al juego: %v", err)
	}

	log.Printf("Jugador: %s", r.GetMessage())

	comienzaeljuego := false

	for !comienzaeljuego {
		time.Sleep(5 * time.Second)
		message := "Ya comienza el juego??"
		r, err := client.BeginGame(ctx, &pb.BeginRequest{Id: id, Message: message})
		if err != nil {
			log.Fatalf("Error en comenzar el juego: %v", err)
		}
		comienzaeljuego = r.GetMessage()
		Jugadoresonline := r.GetOnline()
		fmt.Println(">> Jugadores online: %d", Jugadoresonline)
	}

	// Solicitud Etapa 1
	stage := "¿Va comenzar la etapa?"
	comienzalaetapa := false

	for !comienzalaetapa {

		time.Sleep(10 * time.Second)
		r, err := client.BeginStage(ctx, &pb.BeginStageRequest{Id: id , Stage: stage})
		if err != nil {
			log.Fatalf("Error en comenzar etapa 2: %v", err)
		}


	}






	return
}

func main() {

	// Crear un gRPC canal para comunicarse con el servidor
	// 	-> Esto se crea pasando server address y port number a grpc.Dial()

	// Configurar conexión al server
	conn, err := grpc.Dial(address, grpc.WithInsecure(), grpc.WithBlock())

	if err != nil {
		log.Fatalf("did not connect: %v", err)
	}

	defer conn.Close()

	// Client Stub to perform RPCs
	client := pb.NewGameClient(conn)

	ctx := context.Background()

	for i := 1; i < 16 ; i++ {
		go Bot(ctx, client)
	}

	return

}
