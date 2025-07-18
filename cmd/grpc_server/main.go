package main

import (
	"context"
	"fmt"
	"grpc_with_metrics/internal/interceptor"
	"grpc_with_metrics/internal/metric"
	desc "grpc_with_metrics/pkg/note_v1"
	"log"
	"math/rand"
	"net"
	"net/http"
	"time"

	"github.com/brianvoe/gofakeit"
	"github.com/pkg/errors"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"
	"google.golang.org/protobuf/types/known/timestamppb"
)

const grpcPort = 50051

type server struct {
	desc.UnimplementedNoteV1Server
}

// Get ...
func (s *server) Get(ctx context.Context, req *desc.GetRequest) (*desc.GetResponse, error) {
	if req.GetId() == 0 {
		return nil, errors.Errorf("id is empty")
	}

	// rand.Intn(max - min) + min
	time.Sleep(time.Duration(rand.Intn(1000)) * time.Millisecond)

	return &desc.GetResponse{
		Note: &desc.Note{
			Id: req.GetId(),
			Info: &desc.NoteInfo{
				Title:   gofakeit.BeerName(),
				Content: gofakeit.IPv4Address(),
			},
			CreatedAt: timestamppb.New(gofakeit.Date()),
			UpdatedAt: timestamppb.New(gofakeit.Date()),
		},
	}, nil
}

func main() {
	ctx := context.Background()
	err := metric.Init(ctx)
	if err != nil {
		log.Fatalf("failed to initialize metrics: %v", err)
	}

	lis, err := net.Listen("tcp", fmt.Sprintf(":%d", grpcPort))
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}

	s := grpc.NewServer(grpc.UnaryInterceptor(interceptor.MetricsInterceptor))

	reflection.Register(s)
	desc.RegisterNoteV1Server(s, &server{})

	log.Printf("server listening at %v", lis.Addr())

	// Запускаем http сервер в отдельной горутине для сбора метрик
	go func() {
		err = runPrometheus()
		if err != nil {
			log.Fatal(err)
		}
	}()

	if err := s.Serve(lis); err != nil {
		log.Fatalf("failed to serve: %v", err)
	}
}

// Объявляем http сервер для сбора метрик
func runPrometheus() error {
	mux := http.NewServeMux()
	mux.Handle("/metrics", promhttp.Handler())

	prometheusServer := &http.Server{
		Addr:    "localhost:2112",
		Handler: mux,
	}

	log.Printf("Prometheus server is running on %s", prometheusServer.Addr)
	err := prometheusServer.ListenAndServe()
	if err != nil {
		return err
	}
	return nil
}
