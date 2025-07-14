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
	requestCounter        prometheus.Counter // Создаем метрику типа Counter (счетчкик)
	responseCounter       *prometheus.CounterVec
	histogramResponseTime *prometheus.HistogramVec // Создаем метрику типа Histogram
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
		responseCounter: promauto.NewCounterVec(
			prometheus.CounterOpts{
				Namespace: namespace,
				Subsystem: "grpc",
				Name:      appName + "_response_total",
				Help:      "Количество ответов от сервера",
			}, []string{"status", "method"}, // Прокидываем label. status - строка, которая будет success или error. method - строка, которая отображает метод записавший метрику
		),
		histogramResponseTime: promauto.NewHistogramVec(prometheus.HistogramOpts{ // Инициализируем метрику типа Histogram
			Namespace: namespace,
			Subsystem: "grpc",
			Name:      appName + "_histogram_response_time_seconds",
			Help:      "Время ответа от сервера",
			Buckets:   prometheus.ExponentialBuckets(0.0001, 2, 16),
		}, []string{"status"}), // Также проставляем label, чтобы видеть успешный запрос или нет
	}
	return nil
}

// Функция, которая будет увеличивать счетчик
func IncRequestCounter() {
	metrics.requestCounter.Inc()
}

func IncResponseCounter(status string, method string) {
	metrics.responseCounter.WithLabelValues(status, method).Inc()
}

// Функция, которая будет увеличивать метрику типа Histogram
func HistgramResponseTimeObserve(status string, time float64) {
	metrics.histogramResponseTime.WithLabelValues(status).Observe(time)
}
