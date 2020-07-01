package main

import (
	"context"
	"fmt"
	"io"
	"log"
	"strconv"
	"time"

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
	//doUnary(c)

	// doServerStreaming(c)

	// doClientStreaming(c)
	doBidi(c)
}

func doBidi(c greetpb.GreetServiceClient) {
	log.Println("Starting BiDi")
	// create a bunch of requests
	requests := []*greetpb.GreetEveryoneRequest{
		&greetpb.GreetEveryoneRequest{
			Greeting: &greetpb.Greeting{
				FirstName: "John",
			},
		},
		&greetpb.GreetEveryoneRequest{
			Greeting: &greetpb.Greeting{
				FirstName: "Mary",
			},
		},
		&greetpb.GreetEveryoneRequest{
			Greeting: &greetpb.Greeting{
				FirstName: "Kaarthik",
			},
		},
		&greetpb.GreetEveryoneRequest{
			Greeting: &greetpb.Greeting{
				FirstName: "Kousik",
			},
		},
		&greetpb.GreetEveryoneRequest{
			Greeting: &greetpb.Greeting{
				FirstName: "Vimal",
			},
		},
	}
	// create stream
	stream, err := c.GreetEveryone(context.Background())
	if err != nil {
		log.Fatalf("Unable to create stream: %v", err)
		return
	}

	waitc := make(chan struct{})
	// send bunch of messages - goroutine!
	go func() {
		// func to send bunch of messages
		for _, req := range requests {
			log.Println("Sending message ", req)
			stream.Send(req)
			time.Sleep(1000 * time.Millisecond)
		}
		stream.CloseSend()
	}()
	// recv a bunch of messages - goroutine!
	go func() {
		for {
			res, err := stream.Recv()
			if err == io.EOF {
				break
			}
			if err != nil {
				log.Fatalf("Error recv message: %v", err)
				break
			}
			log.Printf("Got response from server: %v\n", res.GetResult())
		}
		close(waitc)
	}()
	// block till done
	// wait channel
	<-waitc
}
func doClientStreaming(c greetpb.GreetServiceClient) {
	fmt.Println("Starting clientstreaming")
	// get a stream
	stream, err := c.LongGreet(context.Background())

	// send 10 messages
	for i := 0; i < 10; i++ {
		err := stream.Send(&greetpb.LongGreetRequest{
			Greeting: &greetpb.Greeting{
				FirstName: "Aswini " + strconv.Itoa(i),
				LastName:  "K",
			},
		})

		if err != nil {
			log.Fatalf("Failed to send message %v: %v", i, err)
		}
		time.Sleep(100 * time.Millisecond)
	}

	// now get the response
	res, err := stream.CloseAndRecv()
	if err != nil {
		log.Fatalf("Unable to recv message: %v", err)
	}

	// got the message
	log.Printf("Got message: %v\n", res.GetResult())
}

func doServerStreaming(c greetpb.GreetServiceClient) {
	fmt.Println("Starting serverstreaming")

	req := &greetpb.GreetManyTimesRequest{
		Greeting: &greetpb.Greeting{
			FirstName: "Aswini",
			LastName:  "K",
		},
	}

	stream, err := c.GreetManyTimes(context.Background(), req)

	if err != nil {
		log.Fatalf("error calling manytimes: %v", err)
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

		log.Printf("Response from server: %v", msg.GetResult())
	}
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
	log.Printf("Greet: %v", res.Result)
}
