package main

import (
	"context"
	"fmt"
	app_frame "github.com/lianmc123/app-frame"
	"github.com/lianmc123/app-frame/examples/tracing/api/message"
	"github.com/lianmc123/app-frame/examples/tracing/api/sub_message"
	"github.com/lianmc123/app-frame/middleware/tracing"
	"github.com/lianmc123/app-frame/transport/grpc"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/exporters/jaeger"
	"go.opentelemetry.io/otel/sdk/resource"
	tracesdk "go.opentelemetry.io/otel/sdk/trace"
	semconv "go.opentelemetry.io/otel/semconv/v1.4.0"
	grpcx "google.golang.org/grpc"
	"log"
	"math/rand"
	"time"
)

const (
	service     = "trace-message"
	environment = "production"
	id          = 1
)

// tracerProvider returns an OpenTelemetry TracerProvider configured to use
// the Jaeger exporter that will send spans to the provided url. The returned
// TracerProvider will also use a Resource configured with all the information
// about the application.
func setTracerProvider(url string) error {
	// Create the Jaeger exporter
	exp, err := jaeger.New(
		jaeger.WithCollectorEndpoint(jaeger.WithEndpoint(url)),
	)
	if err != nil {
		return err
	}
	tp := tracesdk.NewTracerProvider(
		// Set the sampling rate based on the parent span to 100% 设置采样率
		tracesdk.WithSampler(tracesdk.ParentBased(tracesdk.TraceIDRatioBased(1.0))),
		// Always be sure to batch in production.
		tracesdk.WithBatcher(exp),
		// Record information about this application in an Resource.
		/*tracesdk.WithResource(resource.NewWithAttributes(
			semconv.SchemaURL,
			semconv.ServiceNameKey.String(service),
			attribute.String("environment", environment),
			attribute.Int64("ID", id),
		)),*/
		tracesdk.WithResource(resource.NewSchemaless(
			semconv.ServiceNameKey.String(service),
			attribute.String("env", "dev"),
		)),
	)
	otel.SetTracerProvider(tp)
	return nil
}

type ss struct {
	message.UnimplementedMessageServiceServer
}

func (s *ss) GetMessage(ctx context.Context, req *message.GetMessageReq) (*message.GetMessageReply, error) {
	time.Sleep(time.Millisecond * time.Duration(rand.Intn(500)))
	conn, err := grpc.DialInsecure(ctx,
		grpc.WithEndpoint(":6890"),
		grpc.WithMiddleware(tracing.Client()),
		grpc.WithOptions(grpcx.WithStatsHandler(&tracing.ClientHandler{})),
	)
	if err != nil {
		return nil, err
	}
	client := sub_message.NewGetSubMessageServiceClient(conn)
	resp, err := client.GetSubMessage(ctx, &sub_message.GetSubMessageReq{SubMessage: "hello"})
	if err != nil {
		return nil, err
	}

	reply := message.GetMessageReply{Messages: make([]*message.GetMessageReply_Message, 0, req.Count)}
	for i := 0; i < int(req.Count); i++ {
		reply.Messages = append(reply.Messages, &message.GetMessageReply_Message{Content: fmt.Sprintf("%d Hello %s", i, resp.MessageReply)})
	}

	time.Sleep(time.Millisecond * time.Duration(rand.Intn(500)+ 500))
	return &reply, nil
}

func main() {
	err := setTracerProvider("http://localhost:14268/api/traces")
	if err != nil {
		log.Fatal(err)
	}

	grpcSrv := grpc.NewServer(
		grpc.Address(":6889"),
		grpc.Middleware(tracing.Server()),
	)

	message.RegisterMessageServiceServer(grpcSrv.Server, new(ss))

	app := app_frame.New(
		app_frame.Name("user"),
		app_frame.Service(grpcSrv))
	if err := app.Run(); err != nil {
		log.Fatal(err)
	}
}
