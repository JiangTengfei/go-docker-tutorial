package main

import (
	"context"
	"fmt"
	GoDockerTutorialAdapter "github.com/jiangtengfei/go-docker-tutorial-pub/adapter"
	"github.com/jiangtengfei/go-docker-tutorial-pub/grpc"
	"net/http"
)

const (
	defaultName = "JTF"
)

func indexHandler(w http.ResponseWriter, r *http.Request)  {
	ctx := context.Background()
	reply := GoDockerTutorialAdapter.SayHello(ctx, &grpc.HelloRequest{Name:defaultName})
	_, _ = fmt.Fprint(w, reply.Message)
}
func main() {
	http.HandleFunc("/hi", indexHandler)
	_ = http.ListenAndServe(":8080", nil)
}
