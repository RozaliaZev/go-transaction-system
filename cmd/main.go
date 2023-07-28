package main

import (
	"log"
	"net"
	"net/http"

	ts "go-tr-syst/pkg/gen-tr-syst"

	"google.golang.org/grpc"

	client "go-tr-syst/pkg/client"
	server "go-tr-syst/pkg/server"
)

func main() {
	ginRouter := client.NewRouter()
	grpcServer := grpc.NewServer()
	ts.RegisterTransactionServiceServer(grpcServer, server.GRPCServer{})

	// Запускаем оба сервера параллельно
	go func() {
		lis, err := net.Listen("tcp", ":50051")
		if err != nil {
			log.Fatalf("Failed to listen: %v", err)
		}
		log.Printf("gRPC server listening on :50051")
		grpcServer.Serve(lis)
	}()

	go func() {
		log.Fatal(http.ListenAndServe(":8080", ginRouter))
	}()

	// Чтобы серверы работали бесконечно и не завершались сразу после запуска, можно использовать select{}
	select {}
}
