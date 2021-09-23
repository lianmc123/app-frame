module github.com/lianmc123/app-frame/examples

go 1.16

require (
	github.com/gin-gonic/gin v1.7.4
	github.com/lianmc123/app-frame v0.0.0-00010101000000-000000000000
	github.com/prometheus/client_golang v1.11.0
	go.opentelemetry.io/otel v1.0.0-RC3
	go.opentelemetry.io/otel/exporters/jaeger v1.0.0-RC3
	go.opentelemetry.io/otel/sdk v1.0.0-RC3
	go.uber.org/zap v1.19.0
	google.golang.org/grpc v1.40.0
	google.golang.org/protobuf v1.26.0
)

replace github.com/lianmc123/app-frame => ../
