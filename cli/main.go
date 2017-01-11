package main

import (
	"fmt"
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
	c := pb.NewChargerClient(conn)

	// Contact the server and print out its response.
	stream, err := c.ExecCommand(context.Background(), &pb.CommandRequest{Args: os.Args})
	if err != nil {
		log.Fatalf("could not execute: %v", err)
	}
	for {
		resp, err := stream.Recv()
		if err == io.EOF {
			break
		}
		if err != nil {
			log.Fatalf("error on receiving response: %v", err)
		}
		if resp.StreamType == pb.CommandResponse_STDOUT {
			fmt.Fprintln(os.Stdout, resp.Line)
		} else if resp.StreamType == pb.CommandResponse_STDERR {
			fmt.Fprintln(os.Stderr, resp.Line)
		}
	}
}
