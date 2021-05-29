package main

import (
	"context"
	"fmt"
	"io"
	"log"

	"github.com/vvezani/grpc-go-course/greet/greetpb"
	"google.golang.org/grpc"
)

func main() {
	fmt.Println("Hello, I'm a client")
	cc, err := grpc.Dial("localhost:50051", grpc.WithInsecure())

	if err != nil {
		log.Fatalf("Could not connect: %v", err)
	}
	defer cc.Close()

	c := greetpb.NewGreetServiceClient(cc)

	// doUnary(c)
	doStream(c)
}

func doStream(c greetpb.GreetServiceClient) {
	fmt.Println("Doing Streaming Greet RPC")

	rq := &greetpb.GreetManyTimesRequest{
		Greeting: &greetpb.Greeting{
			FirstName: "Vitor",
			LastName:  "Vezani",
		},
	}

	stream, err := c.GreetManyTimes(context.Background(), rq)
	if err != nil {
		log.Fatalf("Error while calling GreetManyTimes RPC: %v", err)
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

		log.Printf("Response from greetManyTimes: %v", msg.Result)
	}
}

func doUnary(c greetpb.GreetServiceClient) {
	fmt.Println("Doing Unary Greet RPC")
	rq := &greetpb.GreetRequest{
		Greeting: &greetpb.Greeting{
			FirstName: "Vitor",
			LastName:  "Vezani",
		},
	}
	res, err := c.Greet(context.Background(), rq)
	if err != nil {
		log.Fatalf("Error while calling greet RPC: %v", err)
	}

	log.Printf("Response from Greet: %v", res.Result)
}
