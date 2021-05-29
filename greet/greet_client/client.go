package main

import (
	"context"
	"fmt"
	"io"
	"log"
	"time"

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
	// doServerStream(c)
	// doClientStream(c)
	doBiDiStream(c)
}

func doBiDiStream(c greetpb.GreetServiceClient) {
	fmt.Println("Doing BiDi Streaming Greet RPC")

	stream, err := c.GreetEveryone(context.Background())
	if err != nil {
		log.Fatalf("Error while creating stream: %v", err)
	}

	requests := []*greetpb.GreetEveryoneRequest{
		{
			Greeting: &greetpb.Greeting{
				FirstName: "Vitor",
			},
		},
		{
			Greeting: &greetpb.Greeting{
				FirstName: "Lucy",
			},
		},
		{
			Greeting: &greetpb.Greeting{
				FirstName: "Stephan",
			},
		},
		{
			Greeting: &greetpb.Greeting{
				FirstName: "Mark",
			},
		},
	}

	waitc := make(chan struct{})

	go func() {
		for _, req := range requests {
			fmt.Printf("Sending message: %v\n", req)
			stream.Send(req)
			time.Sleep(1000 * time.Millisecond)
		}
		stream.CloseSend()
	}()

	go func() {
		fmt.Println("Receiving messages from server")
		for {
			res, err := stream.Recv()
			if err == io.EOF {
				log.Println("Server finished sending info")
				break
			}
			if err != nil {
				log.Fatalf("Error while receiving GreetEveryoneResponse RPC: %v\n", err)
				break
			}
			fmt.Printf("Received: %v\n", res.GetResult())
		}
		close(waitc)
	}()

	fmt.Println("Client blocked")
	<-waitc
	fmt.Println("Finishing client...")
}

func doClientStream(c greetpb.GreetServiceClient) {
	fmt.Println("Doing Client Streaming Greet RPC")

	stream, err := c.LongGreet(context.Background())
	if err != nil {
		log.Fatalf("Error while calling LongGreet RPC: %v", err)
	}

	requests := []*greetpb.LongGreetRequest{
		{
			Greeting: &greetpb.Greeting{
				FirstName: "Vitor",
			},
		},
		{
			Greeting: &greetpb.Greeting{
				FirstName: "Lucy",
			},
		},
		{
			Greeting: &greetpb.Greeting{
				FirstName: "Stephan",
			},
		},
		{
			Greeting: &greetpb.Greeting{
				FirstName: "Mark",
			},
		},
	}

	for _, rq := range requests {
		stream.Send(rq)
		time.Sleep(100 * time.Millisecond)
	}
	res, err := stream.CloseAndRecv()
	if err != nil {
		log.Fatalf("Error while receiving LongGreetResponse RPC: %v", err)
	}

	log.Printf("Response from LongGreet: %v", res.Result)
}

func doServerStream(c greetpb.GreetServiceClient) {
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

		log.Printf("Response from GreetManyTimes: %v", msg.Result)
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
