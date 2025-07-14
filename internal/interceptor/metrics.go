package interceptor

import (
	"context"
	"grpc_with_metrics/internal/metric"

	"google.golang.org/grpc"
)

// Функция интерцептора, которая перед каждым запросом будет увеличивать счетчик
func MetricsInterceptor(ctx context.Context, req interface{}, info *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (resp interface{}, err error) {
	metric.IncRequestCounter()

	return handler(ctx, req)
}
