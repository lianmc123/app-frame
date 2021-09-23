package app_frame

import (
	"context"
	"fmt"
	uuid "github.com/satori/go.uuid"
	"testing"
	"time"
)

type TestSvc struct {
	done chan struct{}
}

func (t *TestSvc) Start(ctx context.Context) error {
	if t.done == nil {
		t.done = make(chan struct{}, 1)
	}
	for {
		select {
		case <-t.done:
			return nil
		default:
			fmt.Println(uuid.NewV4().String())
			time.Sleep(time.Second)
		}
	}
}

func (t *TestSvc) Stop(ctx context.Context) error {
	if t.done != nil {
		t.done <- struct{}{}
		fmt.Println("STOP")
	}
	return nil
}

func TestApp(t *testing.T) {
	app := New(
		Name("app-frame"),
		Version("v1.0.0"),
		Service(&TestSvc{done: make(chan struct{}, 1)}),
	)
	time.AfterFunc(time.Second*5, func() {
		app.Stop()
	})
	if err := app.Run(); err != nil {
		t.Fatal(err)
	}
}
