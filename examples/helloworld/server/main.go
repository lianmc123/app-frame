package main

import (
	"context"
	"errors"
	"fmt"
	app_frame "github.com/lianmc123/app-frame"
	helloworld "github.com/lianmc123/app-frame/examples/helloworld/proto"
	"github.com/lianmc123/app-frame/transport/grpc"
	"log"
)

type server struct {
	helloworld.UnimplementedGreeterServer
}

func (s *server) SayHello(ctx context.Context, req *helloworld.HelloRequest) (*helloworld.HelloReply, error) {
	if req.Name == "error" {
		return nil, errors.New("fuck error")
	}
	if req.Name == "panic" {
		panic("transport panic")
	}
	fmt.Println("hello")
	return &helloworld.HelloReply{Message: req.Name + " Hello"}, nil
}

func main() {
	srv := grpc.NewServer(grpc.Address(":7892"))
	if ep, err := srv.Endpoint(); err == nil {
		fmt.Println(ep.Host)
	}
	helloworld.RegisterGreeterServer(srv.Server, &server{})

	app := app_frame.New(
		app_frame.Name("helloworld_srv"),
		app_frame.Service(srv),
	)
	if err := app.Run(); err != nil {
		log.Fatal(err)
	}
}
