package main

import (
	"context"
	"fmt"
	"io"
	"log"
	"net"

	"github.com/vvezani/grpc-go-course/calculator/calculatorpb"

	"google.golang.org/grpc"
)

type server struct {
	calculatorpb.UnimplementedCalculatorServiceServer
}

func (*server) ComputeAverage(stream calculatorpb.CalculatorService_ComputeAverageServer) error {
	fmt.Println("ComputeAverage invoked with greeting")
	var sum = int32(0)
	var count = 0
	for {
		msg, err := stream.Recv()
		if err == io.EOF {
			var average = float64(sum) / float64(count)
			stream.SendAndClose(&calculatorpb.ComputeAverageResponse{
				Average: average,
			})
			break
		}
		if err != nil {
			log.Fatalf("Erro while reading client stream: %v", err)
		}

		sum += msg.GetNumber()
		count++
	}
	return nil
}

func (*server) Sum(ctx context.Context, rq *calculatorpb.SumRequest) (*calculatorpb.SumResponse, error) {
	fmt.Printf("Calculator invoked with calculatoring: %v", rq)
	firstNumber := rq.GetFirstNumber()
	secondNumber := rq.GetSecondNumber()

	sum := firstNumber + secondNumber

	res := &calculatorpb.SumResponse{
		SumResult: sum,
	}

	return res, nil
}

func (*server) PrimeNumberDecomposition(rq *calculatorpb.PrimeNumberDecompositionRequest, stream calculatorpb.CalculatorService_PrimeNumberDecompositionServer) error {
	fmt.Printf("PrimeNumberDecomposition invoked with rq: %v", rq)
	number := rq.GetNumber()
	divisor := int64(2)
	for number > 1 {
		if number%divisor == 0 {
			stream.Send(
				&calculatorpb.PrimeNumberDecompositionResponse{
					PrimeFactor: divisor,
				})
			number = number / divisor
		} else {
			divisor++
			fmt.Printf("Divisor has increased to %v\n", divisor)
		}
	}
	return nil
}

func main() {
	fmt.Println("Calculator Service")

	lis, err := net.Listen("tcp", "0.0.0.0:50051")
	if err != nil {
		log.Fatalf("Failed to listen: %v", err)
	}

	s := grpc.NewServer()
	calculatorpb.RegisterCalculatorServiceServer(s, &server{})

	if err := s.Serve(lis); err != nil {
		log.Fatalf("Failed to serve: %v", err)
	}
}
