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
	"os"
	"io"
)

const (
	port = ":9000"
)

type server struct{}

type streamWriter struct {
	stream pb.Charger_ExecCommandServer
	streamType pb.CommandResponse_StreamType
}

var (
	app *cli.App
)

func (s *server) ExecCommand(in *pb.CommandRequest, stream pb.Charger_ExecCommandServer) error {
	stdout := &streamWriter{stream: stream, streamType: pb.CommandResponse_STDOUT}
	stderr := &streamWriter{stream: stream, streamType: pb.CommandResponse_STDERR}

	// Set up CLI handler
	app = cli.NewApp()
	app.Usage = os.Getenv("USAGE")
	app.Version = os.Getenv("VERSION")
	if app.Version == "" {
		app.HideVersion = true
	}

	app.Writer = stdout
	app.ErrWriter = stderr
	app.Name = filepath.Base(in.Args[0])
	app.UsageText = fmt.Sprintf("%s COMMAND", filepath.Base(in.Args[0]))
	err := app.Run(in.Args)
	stream.Send(&pb.CommandResponse{ExitCode: int32(HandleExitCoder(err, stderr))})
	return nil
}

func HandleExitCoder(err error, w io.Writer) int {
	if err == nil {
		return 0
	}

	if exitErr, ok := err.(cli.ExitCoder); ok {
		if err.Error() != "" {
			if _, ok := exitErr.(cli.ErrorFormatter); ok {
				fmt.Fprintf(w, "%+v\n", err)
			} else {
				fmt.Fprintln(w, err)
			}
		}
		return exitErr.ExitCode()
	}

	if multiErr, ok := err.(cli.MultiError); ok {
		for _, merr := range multiErr.Errors {
			HandleExitCoder(merr, w)
		}
		return 1
	}

	if err.Error() != "" {
		if _, ok := err.(cli.ErrorFormatter); ok {
			fmt.Fprintf(w, "%+v\n", err)
		} else {
			fmt.Fprintln(w, err)
		}
	}
	return 1
}

func (s *streamWriter) Write(p []byte) (n int, err error) {
	if err := s.stream.Send(&pb.CommandResponse{Buffer: p, StreamType: s.streamType}); err != nil {
		return 0, err
	}
	return len(p), nil
}

func main() {
	// Do not crash on error
	cli.OsExiter = func(c int) {
		log.Printf("Command returned exit code %d\n", c)
	}

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
