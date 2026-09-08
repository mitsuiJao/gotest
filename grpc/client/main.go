package main

import (
	"context"
	"io"
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

	// Add()
	res, err := client.Add(context.Background(), &pb.TwoNumbers{A: 3, B: 4})
	if err != nil {
		log.Fatalf("Add failed: %v", err)
	}
	log.Println("Add result: ", res.Value)

	// CountUp()
	stream, err := client.CountUp(context.Background(), &pb.CountUpRequest{From: 1, To: 5})
	if err != nil {
		log.Fatalf("CountUp failed: %v", err)
	}

	for {
		res, err := stream.Recv()
		if err == io.EOF {
			break
		}
		if err != nil {
			log.Fatalf("stream error: %v", err)
		}
		log.Println("CountUp: ", res.Value)
	}

	// Sum()
	sumStream, err := client.Sum(context.Background())
	if err != nil {
		log.Fatalf("Sum failed: %v", err)
	}

	nums := []*pb.TwoNumbers{
		{A: 1, B: 2},
		{A: 3, B: 4},
		{A: 5, B: 6},
	}

	for _, n := range nums {
		if err := sumStream.Send(n); err != nil {
			log.Fatalf("Send failed: %v", err)
		}
	}

	sumRes, err := sumStream.CloseAndRecv()
	if err != nil {
		log.Fatalf("CloseAndRecv falied: %v", err)
	}
	log.Println("Sum result: ", sumRes.Total)
}
