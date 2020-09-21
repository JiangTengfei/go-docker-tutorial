package main

import (
	"context"
	"errors"
	"flag"
	"fmt"
	"github.com/BurntSushi/toml"
	pb "github.com/jiangtengfei/go-docker-tutorial-pub/grpc"
	"go-docker-tutorial/server/config"
	"go.etcd.io/etcd/clientv3"
	"google.golang.org/grpc"
	"log"
	"net"
	"time"
)

type server struct {
}

var appConfig config.AppConfig

func (s *server) SayHello(ctx context.Context, in *pb.HelloRequest) (*pb.HelloReply, error) {
	log.Printf("Received: %v", in.GetName())
	return &pb.HelloReply{Message: "Hello " + in.GetName()}, nil
}


func init() {
	if _, err := toml.DecodeFile("config.toml", &appConfig); err != nil {
		log.Printf("DecodeFile has an error. %+v", err)
	}
	log.Printf("AppConfig: %+v", appConfig)
}

func main() {
	flagArgs, err := parseFlag()
	if err != nil {
		log.Fatalf("parseFlag method has an error: %s", err)
		return
	}

	lis, err := net.Listen("tcp", ":0")
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}

	_, port, err := net.SplitHostPort(lis.Addr().String())

	registService(flagArgs.AppId, ":" + port)

	s := grpc.NewServer()
	pb.RegisterGreeterServer(s, &server{})
	if err := s.Serve(lis); err != nil {
		log.Fatalf("failed to serve: %v", err)
	}
}

func registService(appId, port string) {
	cli, err := clientv3.New(clientv3.Config{
		Endpoints:   []string{"host.docker.internal:2379"},
		DialTimeout: 5 * time.Second,
	})
	if err != nil {
		log.Print("error occurred while create etcd client")
		return
	}
	defer cli.Close()

	ipStr, _ := GetinternalIp()

	ctx, cancel := context.WithTimeout(context.Background(), time.Second*30)
	resp, err := cli.Put(ctx, appId+"/"+ipStr+port, ipStr+port)
	cancel()
	if err != nil {
		fmt.Errorf("error while put: %+v", err)
	}
	log.Printf("resp: %+v", resp)
}

func GetinternalIp() (string, error) {
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

type FlagArgs struct {
	AppId string
}

func parseFlag() (*FlagArgs, error) {
	var appId string
	flag.StringVar(&appId, "appId", "", "AppId")
	flag.Parse()

	if appId == "" {
		return nil, errors.New("AppId is a necessary argument")
	}
	log.Printf("flag args appId: %s", appId)
	return &FlagArgs{AppId: appId,}, nil
}
