package metric

import (
	"context"

	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promauto"
)

const (
	namespace = "my_space"
	appName   = "my_app"
)

// Создаем структуру куда будем определять метрики
type Metrics struct {
	requestCounter prometheus.Counter // Создаем метрику типа Counter (счетчкик)
}

// Создаем глобальный приватный объект, а наружу будут торчать только методы
var metrics *Metrics

// Фнкция инициализации всех метрик
func Init(_ context.Context) error {
	metrics = &Metrics{
		requestCounter: promauto.NewCounter(prometheus.CounterOpts{
			Namespace: namespace,
			Subsystem: "grpc",
			Name:      appName + "_request_total",
			Help:      "Количество запросов к серверу",
		}),
	}
	return nil
}

// Функция, которая будет увеличивать счетчик
func IncRequestCounter() {
	metrics.requestCounter.Inc()
}
