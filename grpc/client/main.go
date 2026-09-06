package main

import (
	"context"
	"log"

	pb "example.com/calc/pb"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

func main() {
	conn, err := grpc.NewClient("localhost:50051", grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		log.Fatalf("connect failed: %v", err)
	}
	defer conn.Close()

	client := pb.NewCalcServiceClient(conn)

	res, err := client.Add(context.Background(), &pb.TwoNumbers{A: 3, B: 4})
	if err != nil {
		log.Fatalf("Add failed: %v", err)
	}
	log.Println("Add result: ", res.Value)
}
