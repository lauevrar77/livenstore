package main

import (
	"fmt"
	"log"
	"net"

	"google.golang.org/grpc"
	pb "livenstore.evrard.online/livenstore_grpc"
	"livenstore.evrard.online/persistance"
	"livenstore.evrard.online/services"
)

func main() {
	es := services.NewEventStore("data", persistance.NewEventWriter, persistance.NewEventReader, persistance.NewStreamWriter)

	lis, err := net.Listen("tcp", fmt.Sprintf("localhost:%d", 5001))
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}
	var opts []grpc.ServerOption
	grpcServer := grpc.NewServer(opts...)
	pb.RegisterLivenstoreServer(grpcServer, &pb.Server{
		ES: es,
	})
	grpcServer.Serve(lis)
}
