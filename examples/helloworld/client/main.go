package main

import (
	"context"
	"fmt"
	helloworld "github.com/lianmc123/app-frame/examples/helloworld/proto"
	"github.com/lianmc123/app-frame/transport/grpc"
	"log"
)

func main() {
	conn, err := grpc.DialInsecure(context.Background(),
		//grpc.WithEndpoint("192.168.56.1:7890"),
		//grpc.WithEndpoint("127.0.0.1:7890"),
		grpc.WithEndpoint(":7892"),
		/*grpc.WithMiddleware(),*/
	)

	if err != nil {
		log.Fatal(err)
	}
	defer conn.Close()
	client := helloworld.NewGreeterClient(conn)
	resp, err := client.SayHello(context.Background(), &helloworld.HelloRequest{Name: "go"})
	if err != nil {
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
