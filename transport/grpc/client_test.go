package grpc

import (
	"context"
	"fmt"
	"testing"
)

func TestDial(t *testing.T) {

}

func TestDialInsecure(t *testing.T) {
	ctx := context.Background()
	conn, err := DialInsecure(ctx, WithEndpoint("grpc://192.168.56.1:8099"))
	if err != nil {
		t.Fatal(err)
	}
	defer conn.Close()
	fmt.Println(conn.GetState())
	fmt.Println(conn.Target())
}
