package main

import (
	"context"
	"fmt"
	"log"
	"net"

	"github.com/vvezani/grpc-go-course/greet/greetpb"

	"google.golang.org/grpc"
)

type server struct {
	greetpb.UnimplementedGreetServiceServer
}

func (*server) Greet(ctx context.Context, rq *greetpb.GreetRequest) (*greetpb.GreetResponse, error) {
	fmt.Printf("Greet invoked with greeting: %v", rq)
	firstName := rq.GetGreeting().GetFirstName()
	lastName := rq.GetGreeting().GetLastName()

	result := "Hello " + firstName + " " + lastName

	res := &greetpb.GreetResponse{
		Result: result,
	}

	return res, nil
}

func main() {
	fmt.Println("Hello World")

	lis, err := net.Listen("tcp", "0.0.0.0:50051")
	if err != nil {
		log.Fatalf("Failed to listen: %v", err)
	}

	s := grpc.NewServer()
	greetpb.RegisterGreetServiceServer(s, &server{})

	if err := s.Serve(lis); err != nil {
		log.Fatalf("Failed to serve: %v", err)
	}
}
