package main

import (
	"io"
	"log"
	"os"

	pb "github.com/fermayo/charger/charger"
	"golang.org/x/net/context"
	"google.golang.org/grpc"
)

const (
	address = "localhost:9000"
)

func main() {
	// Set up a connection to the server.
	conn, err := grpc.Dial(address, grpc.WithInsecure())
	if err != nil {
		log.Fatalf("could not connect: %v", err)
	}
	defer conn.Close()
	c := pb.NewRouterClient(conn)

	// Contact the server and print out its response.
	stream, err := c.ExecCommand(context.Background(), &pb.CommandRequest{Args: os.Args})
	if err != nil {
		log.Fatalf("could not execute: %v", err)
	}
	var exitCode int32 = 0
	for {
		resp, err := stream.Recv()
		if err == io.EOF {
			break
		}
		if err != nil {
			log.Fatalf("error on receiving response: %v", err)
		}
		if resp.StreamType == pb.CommandResponse_STDOUT {
			os.Stdout.Write(resp.Buffer)
		} else if resp.StreamType == pb.CommandResponse_STDERR {
			os.Stderr.Write(resp.Buffer)
		}
		exitCode = resp.ExitCode
	}
	os.Exit(int(exitCode))
}
