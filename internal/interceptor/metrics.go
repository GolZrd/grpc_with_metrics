package interceptor

import (
	"context"
	"grpc_with_metrics/internal/metric"
	"time"

	"google.golang.org/grpc"
)

// Функция интерцептора, которая перед каждым запросом будет увеличивать счетчик
func MetricsInterceptor(ctx context.Context, req interface{}, info *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (resp interface{}, err error) {
	metric.IncRequestCounter()

	// Чтобы наша Histogram метрика работала, нужно передать ей время ответа
	// Поэтому сначала запускаем таймер
	timeStart := time.Now()

	// Интегрируем новую метрику
	res, err := handler(ctx, req)

	// После того, как handler завершился, то сразу считаем разница
	diffTime := time.Since(timeStart)

	// В зависимости от ответа увеличиваем счетчик с тем или иным статусом
	if err != nil {
		metric.IncResponseCounter("error", info.FullMethod)             // Если была ошибка, то увеличиваем счетчик с статусом error
		metric.HistgramResponseTimeObserve("error", diffTime.Seconds()) // В случае неудачи эту метрику увеличиваем с статусом error
	} else {
		metric.IncResponseCounter("success", info.FullMethod)             // Если нет ошибки, то увеличиваем счетчик с статусом success
		metric.HistgramResponseTimeObserve("success", diffTime.Seconds()) // В случае успеха эту метрику увеличиваем с статусом success
	}

	return res, err
}
