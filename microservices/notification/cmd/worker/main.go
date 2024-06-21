package main

import (
	"encoding/json"

	lSdk "github.com/julianojj/essentials-sdk-go/pkg/logger"
	"github.com/julianojj/pdi/notification/internal/adapters"
	"github.com/julianojj/pdi/notification/internal/core/service"
	"github.com/julianojj/pdi/notification/internal/ports"
)

func main() {
	logger := lSdk.NewSlog()
	sqs := adapters.NewSQS()
	notificationRepository := adapters.NewNotificationMongoBD()
	notificationService := service.NewNotificationService(logger, notificationRepository)
	Worker(sqs, notificationService, logger)
	forever := make(chan bool)
	<-forever
}

func Worker(
	queue ports.Queue,
	notificationService *service.NotificationService,
	logger lSdk.Logger,
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
					logger.Error("error to send notification", err)
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
