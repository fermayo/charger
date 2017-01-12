package main

import (
	"log"
	"net"
	"path/filepath"

	"fmt"
	pb "github.com/fermayo/charger/charger"
	"github.com/urfave/cli"
	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"
)

const (
	port = ":9000"
)

type server struct{}

type streamWriter struct {
	stream pb.Charger_ExecCommandServer
}

var (
	app *cli.App
)

func (s *server) ExecCommand(in *pb.CommandRequest, stream pb.Charger_ExecCommandServer) error {
	w := &streamWriter{stream: stream}
	app.Writer = w
	app.ErrWriter = w
	app.Name = filepath.Base(in.Args[0])
	app.UsageText = fmt.Sprintf("%s COMMAND", filepath.Base(in.Args[0]))
	app.Run(in.Args)
	return nil
}

func (s *streamWriter) Write(p []byte) (n int, err error) {
	if err := s.stream.Send(&pb.CommandResponse{Buffer: p}); err != nil {
		return 0, err
	}
	return len(p), nil
}

func main() {
	// Set up CLI handler
	app = cli.NewApp()
	app.Usage = "Charge all the things"

	// Start GRPC server
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
