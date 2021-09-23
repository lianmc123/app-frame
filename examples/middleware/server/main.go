package main

import (
	"context"
	"errors"
	"fmt"
	app_frame "github.com/lianmc123/app-frame"
	helloworld "github.com/lianmc123/app-frame/examples/middleware/proto"
	"github.com/lianmc123/app-frame/middleware"
	"github.com/lianmc123/app-frame/transport/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"log"
)

func aMiddleware(handler middleware.Handler) middleware.Handler {
	return func(ctx context.Context, req interface{}) (reply interface{}, err error) {
		fmt.Println("a middleware in", req)
		reply, err = handler(ctx, req)
		fmt.Println("a middleware out", reply)
		return
	}
}

func bMiddleware(handler middleware.Handler) middleware.Handler {
	return func(ctx context.Context, req interface{}) (reply interface{}, err error) {
		fmt.Println("b middleware in", req)
		reply, err = handler(ctx, req)
		fmt.Println("b middleware out", reply)
		return
	}
}

func cMiddleware(handler middleware.Handler) middleware.Handler {
	return func(ctx context.Context, req interface{}) (reply interface{}, err error) {
		defer func() {
			if recoverErr := recover(); recoverErr != nil {
				err = status.Error(codes.Internal, "Internal panic")
				return
			}
		}()
		reply, err = handler(ctx, req)

		return
	}
}

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
	srv := grpc.NewServer(grpc.Address(":6679"), grpc.Middleware(cMiddleware, aMiddleware, bMiddleware))
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
