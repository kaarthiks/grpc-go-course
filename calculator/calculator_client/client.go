package main

import (
	"context"
	"fmt"
	"io"
	"log"

	"github.com/kaarthiks/grpc-go-course/calculator/calculatorpb"
	"google.golang.org/grpc"
)

func main() {
	fmt.Println("Starting calculator client")
	cc, err := grpc.Dial("localhost:50051", grpc.WithInsecure())

	if err != nil {
		log.Fatalf("Problem opening connection: %v", err)
	}

	c := calculatorpb.NewCalculatorServiceClient(cc)

	doSum(c)

	doPrimeFactorization(c)
}

func doSum(c calculatorpb.CalculatorServiceClient) {
	fmt.Println("Starting to do some sum")
	req := &calculatorpb.SumRequest{
		FirstNumber:  10,
		SecondNumber: 25,
	}

	res, err := c.Sum(context.Background(), req)
	if err != nil {
		log.Fatalf("Failed to send request: %v", err)
	}
	log.Printf("Result of addition: %v", res.SumResult)
}

func doPrimeFactorization(c calculatorpb.CalculatorServiceClient) {
	fmt.Println("Starting prime factorization")
	req := &calculatorpb.PrimeNumberDecompositionRequest{
		Number: 13571947292733,
	}

	// the output of calling the rpc is the stream from which you read the response using Recv()
	stream, err := c.PrimeNumberDecomposition(context.Background(), req)

	if err != nil {
		log.Fatalf("error calling decomposition: %v", err)
	}

	for {
		msg, err := stream.Recv()
		if err == io.EOF {
			// end of stream
			break
		}

		if err != nil {
			log.Fatalf("Error while reading stream: %v", err)
		}

		log.Printf("Prime Factor: %v", msg.GetPrimeFactor())
	}
}
