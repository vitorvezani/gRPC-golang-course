package main

import (
	"context"
	"fmt"
	"log"

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

	doUnary(c)
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
