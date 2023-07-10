package grpc

import (
	"context"
	"google.golang.org/grpc"
	v1 "movieuniverse/pkg/movie/api/proto/v1"
	"log"
	"net"
	"os"
	"os/signal"
)

// RunServer runs gRPC service to publish Movie service
func RunServer(ctx context.Context, v1API v1.MovieServiceServer, port string) error {
	listen, err := net.Listen("tcp", ":"+port)
	if err != nil {
		return err
	}

	// register service
	server := grpc.NewServer()
	v1.RegisterMovieServiceServer(server, v1API)

	// graceful shutdown
	c := make(chan os.Signal, 1)
	signal.Notify(c, os.Interrupt)
	go func() {
		for range c {
			// sig is a ^C, handle it
			log.Println("shutting down gRPC server...")

			server.GracefulStop()

			<-ctx.Done()
		}
	}()

	// start gRPC server
	log.Println("starting gRPC server...")
	return server.Serve(listen)
}
