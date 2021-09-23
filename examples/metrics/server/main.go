package main

import (
	"context"
	"errors"
	"fmt"
	"github.com/gin-gonic/gin"
	app_frame "github.com/lianmc123/app-frame"
	helloworld "github.com/lianmc123/app-frame/examples/metrics/proto"
	"github.com/lianmc123/app-frame/metrics"
	metricsMW "github.com/lianmc123/app-frame/middleware/metrics"
	"github.com/lianmc123/app-frame/transport/grpc"
	"github.com/lianmc123/app-frame/transport/http"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	"log"
	"math/rand"
	stdhttp "net/http"
	"time"
)

var (
	_metricSeconds = prometheus.NewHistogramVec(prometheus.HistogramOpts{
		Namespace: "server",
		Subsystem: "requests",
		Name:      "duration_ms",
		Help:      "server requests duration(ms).",
		Buckets:   []float64{5, 10, 25, 50, 100, 250, 500, 1000},
	}, []string{"kind", "operation"})

	_metricRequests = prometheus.NewCounterVec(prometheus.CounterOpts{
		Namespace: "client",
		Subsystem: "requests",
		Name:      "code_total",
		Help:      "The total number of processed requests",
	}, []string{"kind", "operation", "code", "reason"})
)

type server struct {
	helloworld.UnimplementedGreeterServer
}

func init() {
	prometheus.MustRegister(_metricSeconds, _metricRequests)
}

func (s *server) SayHello(ctx context.Context, req *helloworld.HelloRequest) (*helloworld.HelloReply, error) {
	time.Sleep(time.Millisecond * time.Duration(rand.Intn(900)))
	/*a := rand.Intn(10) % 4
	switch a {
	case 0:
		return nil, apierr.InternalServer("InternalServer", "500")
	case 1:
		return nil, apierr.NotFound("NotFound", "404")
	case 2:
		return nil, apierr.BadRequest("BadRequest", "400")
	}*/
	if req.Name == "error" {
		return nil, errors.New("fuck error")
	}
	if req.Name == "panic" {
		panic("transport panic")
	}
	/*fmt.Println("hello")*/
	return &helloworld.HelloReply{Message: req.Name + " Hello"}, nil
}

func main() {
	srv := grpc.NewServer(
		grpc.Address(":7892"),
		grpc.Middleware(
			metricsMW.Server(
				metricsMW.WithRequests(metrics.NewCounter(_metricRequests)),
				metricsMW.WithSeconds(metrics.NewHistogram(_metricSeconds)),
			),
		),
	)

	httpSrv := http.NewServer(":7893", func(r gin.IRouter) {
		r.GET("/metrics", gin.WrapH(promhttp.Handler()))
		r.GET("/hello/:name", func(c *gin.Context) {
			time.Sleep(time.Millisecond * time.Duration(rand.Intn(900)))
			name := c.Param("name")
			if name == "error" {
				c.JSON(stdhttp.StatusBadRequest, gin.H{"name": "fuck"})
				c.Error(fmt.Errorf("name error fuck"))
				return
			} else if name == "error1" {
				c.JSON(stdhttp.StatusInternalServerError, gin.H{"name": "StatusInternalServerError"})
				//c.Error(fmt.Errorf("%s error", name))
				c.Error(fmt.Errorf("name error ~~~~~"))
				return
			}
			c.JSON(stdhttp.StatusOK, map[string]interface{}{
				"hello": name,
			})
		})
	}, false, http.Middleware(
		metricsMW.Server(
			metricsMW.WithSeconds(metrics.NewHistogram(_metricSeconds)),
			metricsMW.WithRequests(metrics.NewCounter(_metricRequests)),
		),
	))
	/*if ep, err := srv.Endpoint(); err == nil {
		fmt.Println(ep.Host)
	}*/
	helloworld.RegisterGreeterServer(srv.Server, &server{})

	app := app_frame.New(
		app_frame.Name("helloworld_srv"),
		app_frame.Service(srv, httpSrv),
	)
	if err := app.Run(); err != nil {
		log.Fatal(err)
	}
}
