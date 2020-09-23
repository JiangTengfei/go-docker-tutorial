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
	"os"
	"os/signal"
	"syscall"
	"time"
)

var appConfig config.AppConfig
var shutdownChan = make(chan os.Signal)
var unloadChan = make(chan struct{})

type server struct {
}

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

	if err := registService(flagArgs.AppId, ":"+port); err != nil {
		log.Fatalf("failed to regist service: %v", err)
	}

	s := grpc.NewServer()
	pb.RegisterGreeterServer(s, &server{})

	go func() {
		if err := s.Serve(lis); err != nil {
			log.Fatalf("failed to serve: %v", err)
		}
	}()

	signal.Notify(shutdownChan, syscall.SIGINT, syscall.SIGTERM)
	sdSig := <- shutdownChan
	log.Printf("receive shutdown signal: %v", sdSig)
	unRegister()

	time.Sleep(2*time.Second)
	s.GracefulStop()
}

func registService(appId, port string) error {
	cli, err := clientv3.New(clientv3.Config{
		Endpoints:   []string{"host.docker.internal:2379"},
		DialTimeout: 5 * time.Second,
	})
	if err != nil {
		log.Print("error occurred while create etcd client")
		return errors.New("error occurred while create etcd client")
	}

	ipStr, _ := GetinternalIp()

	key := fmt.Sprintf("%s/%s%s", appId, ipStr, port)

	ctx, cancel := context.WithTimeout(context.Background(), time.Second*30)
	defer cancel()

	leaseResp, err := cli.Grant(ctx, 10)
	if err != nil {
		log.Printf("Grant method return err: %v", err)
		return errors.New("Grant method return err: " + err.Error())
	}

	resp, err := cli.Put(ctx, key, ipStr+port, clientv3.WithLease(leaseResp.ID))
	if _, err := cli.KeepAlive(context.TODO(), leaseResp.ID); err != nil {
		return fmt.Errorf("KeepAlive failed. %s", err.Error())
	}

	go func() {
		sig := <- unloadChan
		log.Printf("receive unloadChan signal: %v", sig)
		cli.Delete(context.Background(), key)
	}()

	log.Printf("registService resp: %+v", resp)
	return nil
}

func unRegister()  {
	unloadChan <- struct{}{}
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
