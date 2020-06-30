package main

import (
	"context"
	"fmt"
	"log"

	"github.com/kaarthiks/grpc-go-course/greet/greetpb"
	"google.golang.org/grpc"
)

func main() {
	fmt.Println("Hello from client")
	cc, err := grpc.Dial("localhost:50051", grpc.WithInsecure())

	if err != nil {
		log.Fatalf("Could not connect: %v", err)
	}
	defer cc.Close()

	c := greetpb.NewGreetServiceClient(cc)
	//	fmt.Printf("Created client: %f", c)
	doUnary(c)
}

func doUnary(c greetpb.GreetServiceClient) {
	fmt.Println("Starting odUnary")
	req := &greetpb.GreetRequest{
		Greeting: &greetpb.Greeting{
			FirstName: "Kaarthik",
			LastName:  "S",
		},
	}
	res, err := c.Greet(context.Background(), req)
	if err != nil {
		log.Fatalf("Error calling Greet: %v", err)
	}
	log.Println("Greet: %v", res.Result)
}
