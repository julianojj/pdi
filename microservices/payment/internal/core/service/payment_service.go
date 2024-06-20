package service

import (
	"encoding/base64"
	"encoding/json"
	"fmt"
	"pdi/payment/internal/core/domain"
	"pdi/payment/internal/ports"

	lSdk "github.com/julianojj/essentials-sdk-go/pkg/logger"
	qSdk "github.com/julianojj/essentials-sdk-go/pkg/queue"
)

type (
	PaymentService struct {
		paymentGateway    ports.PaymentGateway
		paymentRepository ports.PaymentRepository
		queue             qSdk.Queue
		logger            lSdk.Logger
	}
	PaymentServiceInput struct {
		Customer PaymentServiceCustomerInput `json:"customer"`
		Order    PaymentServiceOrderInput    `json:"order"`
		Payment  string                      `json:"payment"`
	}
	PaymentServiceCustomerInput struct {
		ID    string `json:"id"`
		Name  string `json:"name"`
		Email string `json:"email"`
	}
	PaymentServiceOrderInput struct {
		ID    string  `json:"order_id"`
		Total float64 `json:"total"`
	}
)

func NewPaymentService(
	paymentGateway ports.PaymentGateway,
	paymentRepository ports.PaymentRepository,
	queue qSdk.Queue,
	logger lSdk.Logger,
) *PaymentService {
	return &PaymentService{
		paymentGateway,
		paymentRepository,
		queue,
		logger,
	}
}

func (p *PaymentService) ProcessPayment(input *PaymentServiceInput) error {
	b, err := base64.StdEncoding.DecodeString(input.Payment)
	if err != nil {
		return err
	}
	var decriptedPaymentToken map[string]any
	if err := json.Unmarshal(b, &decriptedPaymentToken); err != nil {
		return err
	}
	paymentType := fmt.Sprintf("%v", decriptedPaymentToken["method"])
	payment := domain.NewPayment(
		input.Order.ID,
		input.Customer.ID,
		paymentType,
		input.Order.Total,
	)
	paymentGatewayOutput, err := p.paymentGateway.ProcessPayment(map[string]any{
		"method":        payment.Type,
		"amount":        input.Order.Total,
		"payment_token": input.Payment,
	})
	if err != nil {
		return err
	}
	if paymentGatewayOutput["code"] == "0000" {
		payment.AprovePayment()
	} else {
		payment.RejectPayment()
	}
	if err := p.paymentRepository.Save(payment); err != nil {
		return err
	}
	b, err = json.Marshal(map[string]any{
		"order_id":       payment.OrderID,
		"payment_id":     payment.ID,
		"payment_status": payment.Status,
	})
	if err != nil {
		return err
	}
	if err := p.queue.Publish("https://localhost.localstack.cloud:4566/000000000000/confirmed-payment", string(b)); err != nil {
		return err
	}
	b, err = json.Marshal(map[string]any{
		"order_id":       payment.OrderID,
		"payment_id":     payment.ID,
		"payment_status": payment.Status,
		"customer":       input.Customer,
	})
	if err != nil {
		return err
	}
	if err := p.queue.Publish("https://localhost.localstack.cloud:4566/000000000000/notification", string(b)); err != nil {
		return err
	}
	p.logger.Info(
		"process payment",
		map[string]any{
			"payment_id":     payment.ID,
			"order_id":       payment.OrderID,
			"customer_id":    payment.CustomerID,
			"payment_status": payment.Status,
			"payment_code":   paymentGatewayOutput["code"],
			"total":          payment.Amount,
		},
	)
	return nil
}
