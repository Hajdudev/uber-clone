package main

import (
	"context"
	"log"
	"net"
	"os"
	"os/signal"
	"syscall"

	grpcserver "google.golang.org/grpc"
)

var GrcpAddr = ":9092"

func main() {
	ctx, cancel := context.WithCancel(context.Background())

	defer cancel()

	go func() {
		sigCh := make(chan os.Signal, 1)
		signal.Notify(sigCh, os.Interrupt, syscall.SIGTERM)
		<-sigCh
		cancel()
	}()

	lis, err := net.Listen("tcp", GrcpAddr)
	if err != nil {
		log.Fatalf("Failed to listen: %v", err)
	}
	grpcServer := grpcserver.NewServer()

	log.Printf("Starting gRPC service driver Service on port %s", lis.Addr().String())
	go func() {
		if err := grpcServer.Serve(lis); err != nil {
			log.Printf("failed to server %v:", err)
			cancel()
		}
	}()

	// wait for the shutdown signal
	<-ctx.Done()
	log.Println("Shuting down ...")
	grpcServer.GracefulStop()
}
