package main

import (
	"context"
	"fmt"
	helloworld "github.com/lianmc123/app-frame/examples/middleware/proto"
	"github.com/lianmc123/app-frame/transport/grpc"
	"log"
)

func main() {
	conn, err := grpc.DialInsecure(context.Background(),
		grpc.WithEndpoint(":6679"),
		/*grpc.WithMiddleware(),*/
	)

	if err != nil {
		fmt.Println("...")
		log.Fatal(err)
	}
	defer conn.Close()
	client := helloworld.NewGreeterClient(conn)
	resp, err := client.SayHello(context.Background(), &helloworld.HelloRequest{Name: "go"})
	if err != nil {
		fmt.Println("...1")
		log.Fatal(err)
	}
	fmt.Println(resp.Message)

	/*resp2, err := client.SayHello(context.Background(), &helloworld.HelloRequest{Name: "error"})
	if err != nil {
		s, ok := status.FromError(err)
		if ok {
			fmt.Println(s.Code(), s.Message())
		}
		log.Fatal(err)
	}

	fmt.Println(resp2.Message)*/

}
