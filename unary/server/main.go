package main

import (
	"log"
	"main/unary/server/ecommerce"
	"net"

	"google.golang.org/grpc"
)

const port = ":12345"

func main() {
	listen, err := net.Listen("tcp", port)
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}

	s := grpc.NewServer()
	ecommerce.RegisterOrderManagementServer(s, &ecommerce.Server{})

	log.Println("start listen")
	if err := s.Serve(listen); err != nil {
		log.Fatalf("failed to serve: %v", err)
	}
}
