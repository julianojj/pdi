package main

import (
	"encoding/json"
	"pdi/order/internal/adapters"
	"pdi/order/internal/core/service"
	"pdi/order/internal/ports"

	lSdk "github.com/julianojj/essentials-sdk-go/pkg/logger"
)

func main() {
	orderRepository := adapters.NewOrderRepositoryDynamoDB()
	itemRepository := adapters.NewItemRepositoryMemory()
	sqs := adapters.NewSQS()
	userGateway := adapters.NewUserGatewayAPI()
	logger := lSdk.NewSlog()

	orderService := service.NewOrderService(orderRepository, itemRepository, sqs, userGateway, logger)

	Worker(sqs, orderService, logger)
	forever := make(chan bool)
	<-forever
}

func Worker(
	queue ports.Queue,
	orderService *service.OrderService,
	logger lSdk.Logger,
) {
	jobs := []struct {
		name string
		url  string
		fn   func(args []byte) error
	}{
		{
			name: "consumer-confirm-order",
			url:  "https://localhost.localstack.cloud:4566/000000000000/confirmed-payment",
			fn: func(args []byte) error {
				var input map[string]any
				if err := json.Unmarshal(args, &input); err != nil {
					return err
				}
				if err := orderService.UpdateStatusOrder(input); err != nil {
					logger.Error("error to confirm order", err)
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
