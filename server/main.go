package main

import (
	"log"
	"net"
	"fmt"
	"path/filepath"

	pb "github.com/fermayo/charger/charger"
	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"
)

const (
	port = ":9000"
)

type server struct{}

var (
	commands []string;
)

func (s *server) ExecCommand(in *pb.CommandRequest, stream pb.Charger_ExecCommandServer) error {
	if len(in.Args) == 1 {
		if err := stream.Send(&pb.CommandResponse{Line: fmt.Sprintf("\nUsage: %s COMMAND\n", filepath.Base(in.Args[0]))}); err != nil {
			return err
		}
	}
	if len(commands) == 0 {
		if err := stream.Send(&pb.CommandResponse{Line: "No commands registered"}); err != nil {
			return err
		}
	} else {
		if err := stream.Send(&pb.CommandResponse{Line: "Available commands:"}); err != nil {
			return err
		}
		for _, command := range commands {
			if err := stream.Send(&pb.CommandResponse{Line: fmt.Sprintf("   %s", command)}); err != nil {
				return err
			}
		}
	}
	return nil
}

func main() {
	lis, err := net.Listen("tcp", port)
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}
	s := grpc.NewServer()
	pb.RegisterChargerServer(s, &server{})
	reflection.Register(s)
	if err := s.Serve(lis); err != nil {
		log.Fatalf("failed to serve: %v", err)
	}
}
