package interceptor

import (
	"context"
	"grpc_with_metrics/internal/metric"

	"google.golang.org/grpc"
)

// Функция интерцептора, которая перед каждым запросом будет увеличивать счетчик
func MetricsInterceptor(ctx context.Context, req interface{}, info *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (resp interface{}, err error) {
	metric.IncRequestCounter()

	// Интегрируем новую метрику
	res, err := handler(ctx, req)
	// В зависимости от ответа увеличиваем счетчик с тем или иным статусом
	if err != nil {
		metric.IncResponseCounter("error", info.FullMethod) // Если была ошибка, то увеличиваем счетчик с статусом error
	} else {
		metric.IncResponseCounter("success", info.FullMethod) // Если нет ошибки, то увеличиваем счетчик с статусом success
	}

	return res, err
}
