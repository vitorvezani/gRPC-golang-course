package main

import (
	"context"
	"fmt"
	"io"
	"log"
	"time"

	"github.com/vvezani/grpc-go-course/calculator/calculatorpb"
	"google.golang.org/grpc"
)

func main() {
	fmt.Println("Hello, I'm a client")
	cc, err := grpc.Dial("localhost:50051", grpc.WithInsecure())

	if err != nil {
		log.Fatalf("Could not connect: %v", err)
	}
	defer cc.Close()

	c := calculatorpb.NewCalculatorServiceClient(cc)

	// doUnary(c)
	// doStream(c)
	doClientStream(c)
}

func doClientStream(c calculatorpb.CalculatorServiceClient) {
	fmt.Println("Doing Client Streaming Greet RPC")

	stream, err := c.ComputeAverage(context.Background())
	if err != nil {
		log.Fatalf("Error while calling LongGreet RPC: %v", err)
	}

	requests := []*calculatorpb.ComputeAverageRequest{
		{
			Number: 5,
		},
		{
			Number: 15,
		},
		{
			Number: 25,
		},
		{
			Number: 35,
		},
	}

	for _, rq := range requests {
		stream.Send(rq)
		time.Sleep(100 * time.Millisecond)
	}
	res, err := stream.CloseAndRecv()
	if err != nil {
		log.Fatalf("Error while receiving ComputeAverage RPC: %v", err)
	}

	log.Printf("Response from ComputeAverage: %v", res.Average)
}

func doStream(c calculatorpb.CalculatorServiceClient) {
	fmt.Println("Doing Stream Calculator RPC")
	rq := &calculatorpb.PrimeNumberDecompositionRequest{
		Number: 123912932,
	}
	stream, err := c.PrimeNumberDecomposition(context.Background(), rq)
	if err != nil {
		log.Fatalf("Error while calling primeNumberDecomposition RPC: %v", err)
	}

	for {
		msg, err := stream.Recv()
		if err == io.EOF {
			// end of the stream
			log.Println("End of the stream")
			break
		}
		if err != nil {
			log.Fatalf("Error while reading stream: %v", err)
		}

		log.Printf("Response from primeNumberDecomposition: %v", msg.PrimeFactor)
	}
}

func doUnary(c calculatorpb.CalculatorServiceClient) {
	fmt.Println("Doing Unary Calculator RPC")
	rq := &calculatorpb.SumRequest{
		FirstNumber:  5,
		SecondNumber: 45,
	}
	res, err := c.Sum(context.Background(), rq)
	if err != nil {
		log.Fatalf("Error while calling sum RPC: %v", err)
	}

	log.Printf("Response from Calculator: %v", res.SumResult)
}
