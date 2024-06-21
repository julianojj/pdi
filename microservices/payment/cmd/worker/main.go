package main

import (
	"encoding/json"
	"pdi/payment/internal/adapters"
	"pdi/payment/internal/core/service"
	"pdi/payment/internal/ports"

	lSdk "github.com/julianojj/essentials-sdk-go/pkg/logger"
)

func main() {
	sqs := adapters.NewSQS()
	paymentGateway := adapters.NewPagarme()
	paymentRepository := adapters.NewPaymentRepositoryDynamoDB()

	logger := lSdk.NewSlog()

	paymentService := service.NewPaymentService(paymentGateway, paymentRepository, sqs, logger)
	Worker(sqs, paymentService, logger)
	forever := make(chan bool)
	<-forever
}

func Worker(
	queue ports.Queue,
	paymentService *service.PaymentService,
	logger lSdk.Logger,
) {
	jobs := []struct {
		name string
		url  string
		fn   func(args []byte) error
	}{
		{
			name: "consumer-process-payment",
			url:  "https://localhost.localstack.cloud:4566/000000000000/maked-order",
			fn: func(args []byte) error {
				var input *service.PaymentServiceInput
				if err := json.Unmarshal(args, &input); err != nil {
					return err
				}
				if err := paymentService.ProcessPayment(input); err != nil {
					logger.Error("error to process payment", err)
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
