package main

import (
	"context"
	"log"
	"net"

	"github.com/kaarthiks/grpc-go-course/calculator/calculatorpb"
	"google.golang.org/grpc"
)

type server struct{}

func (*server) Sum(ctx context.Context, req *calculatorpb.SumRequest) (*calculatorpb.SumResponse, error) {
	log.Printf("Called with %v\n", req)

	resp := calculatorpb.SumResponse{
		SumResult: req.GetFirstNumber() + req.GetSecondNumber(),
	}
	return &resp, nil
}

func main() {
	log.Println("Starting server")

	lis, err := net.Listen("tcp", "0.0.0.0:50051")
	if err != nil {
		log.Fatalf("Failed to listen: %v", err)
	}

	grpcServer := grpc.NewServer()
	calculatorpb.RegisterCalculatorServiceServer(grpcServer, &server{})

	if err := grpcServer.Serve(lis); err != nil {
		log.Fatalf("Failed to serve: %v", err)
	}
}
