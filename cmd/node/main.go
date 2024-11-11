package main

import (
	"log"
	"net"

	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"
)

func main() {
	listener, err := net.Listen("tcp", ":8080")
	if err != nil {
		log.Println("failed to listen:", err)
	}
	defer listener.Close()

	log.Println("server is listening on", listener.Addr().String())

	s := grpc.NewServer()
	reflection.Register(s)

	if err := s.Serve(listener); err != nil {
		log.Println("failed to serve:", err)
	}
}
