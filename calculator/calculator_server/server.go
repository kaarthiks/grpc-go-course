package main

import (
	"context"
	"io"
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

// The argument to the rpc is the request and the stream.
func (*server) PrimeNumberDecomposition(req *calculatorpb.PrimeNumberDecompositionRequest, stream calculatorpb.CalculatorService_PrimeNumberDecompositionServer) error {
	log.Printf("Received request: %v\n", req)

	number := req.GetNumber()
	//  now run the algorithm to find the prime factors
	// start with divisior as 2. keep dividing by divisor as long as remainder is 0.
	// once remainder is no longer 0, incremeent the divisor and repeat the division

	divisor := uint64(2)

	for number > 1 {
		// log.Printf("Number: %v \t Divisor: %d\n", number, divisor)
		if number%divisor == 0 {
			// divisor is a factor of the number, send it out
			stream.Send(&calculatorpb.PrimeNumberDecompositionResponse{
				PrimeFactor: divisor,
			})
			// divide the number by divisor and repeat cycle
			number /= divisor
		} else {
			// divisor is not a factor of the number, increment the divisor
			divisor++
		}
	}
	// all done
	return nil
}

func (*server) CalcAverage(stream calculatorpb.CalculatorService_CalcAverageServer) error {
	log.Println("Entered calcaverage")

	sum := uint32(0)
	count := 0

	// start the stream and keep going till EOF
	for {
		req, err := stream.Recv()
		if err == io.EOF {
			average := float32(sum) / float32(count)
			// end of numbers
			log.Println("Sending back average: ", average)
			return stream.SendAndClose(&calculatorpb.CalcAverageResponse{
				Average: average,
			})
		}

		if err != nil {
			log.Fatalf("Error reading from stream: %v", err)
		}

		// just another message with a number, compute the sum
		number := req.GetNumber()
		log.Println("Got Number ", number)
		sum += number
		count++
	}
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
