package grpc

import (
	"context"
	"fmt"
	"log"
	"testing"
)

func TestNewServer(t *testing.T) {
	ctx := context.Background()

	srv := NewServer(
		Address(":8099"),
	)

	ep, err := srv.Endpoint()
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println(ep.String())
	if err := srv.Start(ctx); err != nil {
		panic(err)
	}
}


