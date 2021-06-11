package main

import (
	"context"
	"fmt"
	"log"

	"github.com/vvezani/grpc-go-course/blog/blogpb"
	"google.golang.org/grpc"
)

func main() {
	fmt.Println("Blog Client Started")
	cc, err := grpc.Dial("localhost:50051", grpc.WithInsecure())

	if err != nil {
		log.Fatalf("Could not connect: %v", err)
	}
	defer cc.Close()

	c := blogpb.NewBlogServiceClient(cc)

	fmt.Println("Creating blog request...")
	req := &blogpb.CreateBlogRequest{
		Blog: &blogpb.Blog{
			AuthorId: "Vitor",
			Title:    "My First Blog",
			Content:  "Content of my first blog",
		},
	}

	res, err := c.CreateBlog(context.Background(), req)
	if err != nil {
		log.Fatalf("Unexpected error: %v", err)
	}
	fmt.Printf("Blog has been created: %v", res.Blog)

	// read blog

	fmt.Println("Reading blog 1")

	_, err2 := c.ReadBlog(context.Background(), &blogpb.ReadBlogRequest{
		BlogId: "123",
	})
	if err2 != nil {
		fmt.Printf("Error while reading the blog: %v\n", err2)
	}

	fmt.Println("Reading blog 2")
	res2, err3 := c.ReadBlog(context.Background(), &blogpb.ReadBlogRequest{
		BlogId: res.Blog.Id,
	})
	if err3 != nil {
		fmt.Printf("Error while reading the blog: %v\n", err3)
	}

	fmt.Printf("Blog was read: %v\n", res2)
}
