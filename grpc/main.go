package main

import (
	"context"
	"log"
	"net"

	pb "example.com/calc/pb"
	"google.golang.org/grpc"
)

type server struct {
	pb.UnimplementeCalcServiceServer
}

func (s *server) Add(ctx context.Context, req *pb.TwoNumbers) (*pb.Result, error) {
	return &pb.Result{Value: req.A + req.B}, nil
}
func (s *server) Multiply(ctx context.Context, req *pb.TwoNumbers) (*pb.Result, error) {
	return &pb.Result{Value: req.A * req.B}, nil
}
func (s *server) CountUp(req *pb.CountUpRequest, stream pb.CalcService_CountUpServer) error {
	for i := req.From; i <= req.To; i++ {
		if err := stream.Send(&pb.Result{Value: i}); err != nil {
			return err
		}
	}
	return nil
}

func main() {
	lis, err := net.Listen("tcp", ":50051")
	if err != nil {
		log.Fatal("listen failed: %v", err)
	}

	s := grpc.NewServer()
	pb.RegisterCalcServiceServer(s, &server{})

	log.Println("calc-server listening on :50051")
	if err := s.Serve(lis); err != nil {
		log.Fatal("serve failed: %v", err)
	}
}
