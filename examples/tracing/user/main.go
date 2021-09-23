package main

import (
	"fmt"
	"github.com/gin-gonic/gin"
	app_frame "github.com/lianmc123/app-frame"
	"github.com/lianmc123/app-frame/examples/tracing/api/message"
	"github.com/lianmc123/app-frame/middleware/tracing"
	"github.com/lianmc123/app-frame/transport/grpc"
	"github.com/lianmc123/app-frame/transport/http"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/exporters/jaeger"
	"go.opentelemetry.io/otel/sdk/resource"
	tracesdk "go.opentelemetry.io/otel/sdk/trace"
	semconv "go.opentelemetry.io/otel/semconv/v1.4.0"
	grpcx "google.golang.org/grpc"
	"log"
	"math/rand"
	stdhttp "net/http"
	"time"
)

const (
	service     = "trace-demo"
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

func main() {
	err := setTracerProvider("http://localhost:14268/api/traces")
	if err != nil {
		log.Fatal(err)
	}

	httpSrv := http.NewServer(":6888", func(r gin.IRouter) {
		r.GET("/msg", func(c *gin.Context) {
			conn, err := grpc.DialInsecure(c.Request.Context(),
				grpc.WithEndpoint(":6889"),
				grpc.WithMiddleware(tracing.Client()),
				grpc.WithOptions(grpcx.WithStatsHandler(&tracing.ClientHandler{})),
			)
			time.Sleep(time.Millisecond * time.Duration(rand.Intn(500)))
			if err != nil {
				c.JSON(stdhttp.StatusInternalServerError, map[string]interface{}{
					"msg": "报错了",
				})
				return
			}
			client := message.NewMessageServiceClient(conn)
			resp, err := client.GetMessage(c.Request.Context(), &message.GetMessageReq{
				Id:    1,
				Count: int64(rand.Intn(5)),
			})
			if err != nil {
				c.JSON(stdhttp.StatusInternalServerError, map[string]interface{}{
					"msg": err.Error(),
				})
				return
			}
			for _, replyMessage := range resp.Messages {
				fmt.Println(replyMessage.Content)
			}
			c.JSON(stdhttp.StatusOK, map[string]interface{}{
				"msg": "ok",
			})
			return
		})
	}, true, http.Middleware(
		//tracing.Server(),
	))

	app := app_frame.New(
		app_frame.Name("user"),
		app_frame.Service(httpSrv))
	if err := app.Run(); err != nil {
		log.Fatal(err)
	}
}
