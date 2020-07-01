package main

import (
	"context"
	"fmt"
	"io"
	"log"
	"net"
	"strconv"
	"time"

	"github.com/kaarthiks/grpc-go-course/greet/greetpb"
	"google.golang.org/grpc"
)

type server struct{}

func (*server) Greet(ctx context.Context, req *greetpb.GreetRequest) (*greetpb.GreetResponse, error) {
	fmt.Printf("Greet called with: %v\n", req)
	firstName := req.GetGreeting().GetFirstName()
	result := "Hello " + firstName
	resp := greetpb.GreetResponse{
		Result: result,
	}
	return &resp, nil
}

// Input is the request and the stream
func (*server) GreetManyTimes(req *greetpb.GreetManyTimesRequest, stream greetpb.GreetService_GreetManyTimesServer) error {
	log.Println("Invoked Manytimesserver")
	firstName := req.GetGreeting().GetFirstName()

	for i := 0; i < 10; i++ {
		result := "Hello " + firstName + " number " + strconv.Itoa(i)

		resp := &greetpb.GreetManyTimesResponse{
			Result: result,
		}
		stream.Send(resp)
		time.Sleep(1 * time.Second)
	}
	return nil
}

// the input is the stream
func (*server) LongGreet(stream greetpb.GreetService_LongGreetServer) error {
	// Do nothing for now
	log.Println("In LongGreet now")
	greeting := ""
	for {
		msg, err := stream.Recv()
		if err == io.EOF {
			// received end of messages, send response back
			err := stream.SendAndClose(&greetpb.LongGreetResponse{
				Result: greeting,
			})
			if err != nil {
				log.Fatalf("Got error sending response: %v", err)
			}

			// break from this loop
			break
		}
		// if a real error is received
		if err != nil {
			log.Fatalf("Got error in recv: %v", err)
		}

		// got a message, print it out
		log.Printf("Got message from client: %v", msg)
		greeting = greeting + "Hello " + msg.GetGreeting().GetFirstName() + "!\n"
	}

	return nil
}

func (*server) GreetEveryone(stream greetpb.GreetService_GreetEveryoneServer) error {
	log.Println("Invoked GreetEveryone")

	for {
		req, err := stream.Recv()
		if err == io.EOF {
			return nil
		}
		if err != nil {
			log.Fatalf("Error in recv: %v", err)
			return err
		}
		firstName := req.GetGreeting().GetFirstName()
		result := "Hello " + firstName
		senderr := stream.Send(&greetpb.GreetEveryoneResponse{
			Result: result,
		})
		if senderr != nil {
			log.Fatalf("Error while sending data to client: %v", senderr)
			return senderr
		}
	}
}

func main() {
	fmt.Println("Hello world")

	lis, err := net.Listen("tcp", "0.0.0.0:50051")

	if err != nil {
		log.Fatalf("Failed to listen: %v", err)
	}

	s := grpc.NewServer()
	greetpb.RegisterGreetServiceServer(s, &server{})

	if err := s.Serve(lis); err != nil {
		log.Fatalf("failed to serve: %v", err)
	}
}
