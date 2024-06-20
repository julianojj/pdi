package main

import (
	"encoding/json"
	"log/slog"
	"os"

	"github.com/julianojj/pdi/notification/internal/adapters"
	"github.com/julianojj/pdi/notification/internal/core/service"
	"github.com/julianojj/pdi/notification/internal/ports"
)

func main() {
	logger := slog.New(slog.NewJSONHandler(os.Stdout, nil)).With(
		slog.Any("application", map[string]any{
			"name":        "order-ms",
			"environment": "dev",
			"version":     "1.0.0",
		}),
	)
	sqs := adapters.NewSQS()
	notificationService := service.NewNotificationService(logger)
	Worker(sqs, notificationService, logger)
	forever := make(chan bool)
	<-forever
}

func Worker(
	queue ports.Queue,
	notificationService *service.NotificationService,
	logger *slog.Logger,
) {
	jobs := []struct {
		name string
		url  string
		fn   func(args []byte) error
	}{
		{
			name: "consumer-notification",
			url:  "https://localhost.localstack.cloud:4566/000000000000/notification",
			fn: func(args []byte) error {
				var input map[string]any
				if err := json.Unmarshal(args, &input); err != nil {
					return err
				}
				if err := notificationService.NotifyPaymentOrder(input); err != nil {
					logger.Error(
						"error to send notification",
						slog.Any("data", map[string]any{
							"err": err,
						}),
					)
					return err
				}
				return nil
			},
		},
	}
	for _, job := range jobs {
		go queue.Consume(job.url, job.fn)
	}
}
