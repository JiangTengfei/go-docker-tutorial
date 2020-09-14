package main

import (
	"context"
	"errors"
	"fmt"
	"github.com/coreos/etcd/clientv3"
	pb "github.com/jiangtengfei/go-docker-tutorial-pub/grpc"
	"google.golang.org/grpc"
	"log"
	"net"
	"time"
)

const port = ":9090"

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

	ipStr, _ := GetInterIp()

	ctx, cancel := context.WithTimeout(context.Background(), time.Second*30)
	resp, err := cli.Put(ctx, "server_ip", ipStr+":9090")
	cancel()
	if err != nil {
		fmt.Errorf("error while put: %+v", err)
	}
	log.Printf("resp: %+v", resp)
}

func GetInterIp() (string, error) {
	addrs, err := net.InterfaceAddrs()
	if err != nil {
		return "", err
	}

	for _, address := range addrs {
		// check the address type and if it is not a loopback the display it
		if ipnet, ok := address.(*net.IPNet); ok && !ipnet.IP.IsLoopback() {
			if ipnet.IP.To4() != nil {
				//fmt.Println(ipnet.IP.String())
				return ipnet.IP.String(), nil
			}
		}
	}

	return "", errors.New("no inter ip")
}
