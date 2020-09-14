package main

import (
	"context"
	"fmt"
	"github.com/coreos/etcd/clientv3"
	pb "github.com/jiangtengfei/go-docker-tutorial-pub/grpc"
	"google.golang.org/grpc"
	"log"
	"net"
	"time"
)

const port  = ":9090"

type server struct {
}

func (s *server) SayHello(ctx context.Context, in *pb.HelloRequest) (*pb.HelloReply, error) {
	log.Printf("Received: %v", in.GetName())
	return &pb.HelloReply{Message: "Hello " + in.GetName()}, nil
}

func main() {
	lis, err := net.Listen("tcp", port)
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}

	registService()

	s := grpc.NewServer()
	pb.RegisterGreeterServer(s, &server{})
	if err := s.Serve(lis); err != nil {
		log.Fatalf("failed to serve: %v", err)
	}
}

func registService() {
	cli, err := clientv3.New(clientv3.Config{
		Endpoints:   []string{"127.0.0.1:2379"},
		DialTimeout: 5 * time.Second,
	})
	if err != nil {
		log.Print("error occurred while create etcd client")
		return
	}
	defer cli.Close()

	ctx, cancel := context.WithTimeout(context.Background(), time.Second*30)
	resp, err := cli.Put(ctx, "server_ip", "127.0.0.1:9090")
	cancel()
	if err != nil {
		fmt.Errorf("error while put: %+v", err)
	}
	log.Printf("resp: %+v", resp)
}
