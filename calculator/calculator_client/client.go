package main

import (
	"context"
	"fmt"
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
