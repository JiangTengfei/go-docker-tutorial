package main

import (
	"context"
	"fmt"
	GoDockerTutorialAdapter "github.com/jiangtengfei/go-docker-tutorial-pub/adapter"
	"github.com/jiangtengfei/go-docker-tutorial-pub/grpc"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"
)

const (
	defaultName = "JTF"
)

func indexHandler(w http.ResponseWriter, r *http.Request) {
	ctx := context.Background()
	reply := GoDockerTutorialAdapter.SayHello(ctx, &grpc.HelloRequest{Name: defaultName})
	_, _ = fmt.Fprint(w, reply.Message)
}
func main() {

	//graceful shutdown
	ch := make(chan os.Signal)

	http.HandleFunc("/hi", indexHandler)
	server := &http.Server{Addr: ":8080", Handler: nil}
	go func() {
		_ = server.ListenAndServe()
	}()

	signal.Notify(ch, syscall.SIGINT, syscall.SIGTERM)
	sig := <-ch
	log.Printf("receive syscall: %v", sig)

	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()
	log.Println(server.Shutdown(ctx))
}
