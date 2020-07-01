package main

import (
	"context"
	"fmt"
	"io"
	"log"
	"math/rand"
	"time"

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

	// doSum(c)

	// doPrimeFactorization(c)

	// doCalcAverage(c)

	doMax(c)
}

func doMax(c calculatorpb.CalculatorServiceClient) {
	log.Println("Doing doMax")

	stream, err := c.FindMax(context.Background())
	if err != nil {
		log.Fatalf("Error opening stream: %v", err)
	}
	waitc := make(chan struct{})
	// send go routine
	go func() {
		numbers := []int32{4, 7, 2, 11, 2, 3, 13, 42, 3}
		for _, num := range numbers {
			log.Println("Sending number ", num)
			senderr := stream.Send(&calculatorpb.FindMaxRequest{
				Number: num,
			})
			if senderr != nil {
				log.Fatalf("Error sending number: %v", senderr)
			}
			time.Sleep(1 * time.Second)
		}
		stream.CloseSend()
	}()

	// recv
	go func() {
		for {
			res, err := stream.Recv()
			if err == io.EOF {
				log.Println("Finished all numbers")
				break
			}
			if err != nil {
				log.Fatalf("ERror receiving max: %v", err)
				break
			}
			// got max number, print
			max := res.GetMax()
			log.Println("Max number is ", max)
		}
		close(waitc)
	}()

	// block on waitc
	<-waitc
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

func doCalcAverage(c calculatorpb.CalculatorServiceClient) {
	fmt.Println("Calculating average")

	stream, err := c.CalcAverage(context.Background())
	if err != nil {
		log.Fatalf("Unable to setup stream")
	}

	// start sending numbers
	for i := 0; i < 10; i++ {
		req := &calculatorpb.CalcAverageRequest{
			Number: uint32(rand.Intn(100)),
		}
		log.Printf("Sending number: %v", req)
		err := stream.Send(req)
		if err != nil {
			log.Fatalf("Cannot send number to server: %v", err)
		}
	}

	// now receive
	res, err := stream.CloseAndRecv()
	if err != nil {
		log.Fatalf("Unable to recv from server: %v", err)
	}
	log.Printf("Average: %v", res)
}
