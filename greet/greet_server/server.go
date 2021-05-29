package main

import (
	"context"
	"fmt"
	"io"
	"log"
	"net"
	"strconv"
	"time"

	"github.com/vvezani/grpc-go-course/greet/greetpb"

	"google.golang.org/grpc"
)

type server struct {
	greetpb.UnimplementedGreetServiceServer
}

func (*server) LongGreet(stream greetpb.GreetService_LongGreetServer) error {
	fmt.Println("LongGreet invoked with greeting")
	var result = ""
	for {
		msg, err := stream.Recv()
		if err == io.EOF {
			stream.SendAndClose(&greetpb.LongGreetResponse{
				Result: result,
			})
			break
		}
		if err != nil {
			log.Fatalf("Erro while reading client stream: %v", err)
		}

		var firstName = msg.GetGreeting().GetFirstName()
		result += firstName + "! "
	}
	return nil
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

func (*server) GreetManyTimes(rq *greetpb.GreetManyTimesRequest, stream greetpb.GreetService_GreetManyTimesServer) error {
	fmt.Printf("GreetManyTimes invoked with greeting: %v", rq)
	firstName := rq.GetGreeting().GetFirstName()

	for i := 0; i < 10; i++ {
		result := "Hello " + firstName + " number " + strconv.Itoa(i)
		res := &greetpb.GreetManyTimesResponse{
			Result: result,
		}
		stream.Send(res)
		time.Sleep(1000 * time.Millisecond)
	}
	return nil
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
