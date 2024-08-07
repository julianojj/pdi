package service

import (
	"encoding/base64"
	"encoding/json"
	"fmt"
	"pdi/order/internal/core/domain"
	"pdi/order/internal/core/exceptions"
	"pdi/order/internal/ports"

	lSdk "github.com/julianojj/essentials-sdk-go/pkg/logger"
	qSdk "github.com/julianojj/essentials-sdk-go/pkg/queue"
)

type (
	OrderService struct {
		orderRepository ports.OrderRepository
		itemRepository  ports.ItemRepository
		queue           qSdk.Queue
		userGateway     ports.UserGateway
		logger          lSdk.Logger
	}
	OrderServiceInput struct {
		UserID       string                   `json:"user_id"`
		PaymentToken string                   `json:"payment_token"`
		Items        []*OrderServiceItemInput `json:"items"`
	}
	OrderServiceItemInput struct {
		ID       string `json:"item_id"`
		Quantity int    `json:"quantity"`
	}
)

func NewOrderService(
	orderRepository ports.OrderRepository,
	itemRepository ports.ItemRepository,
	queue qSdk.Queue,
	userGateway ports.UserGateway,
	logger lSdk.Logger,
) *OrderService {
	return &OrderService{
		orderRepository,
		itemRepository,
		queue,
		userGateway,
		logger,
	}
}

func (os *OrderService) MakeOrder(input *OrderServiceInput) (map[string]any, error) {
	if input.UserID == "" {
		return nil, exceptions.ErrUserIDIsRequired
	}
	if input.PaymentToken == "" {
		return nil, exceptions.ErrPaymentTokenIsRequired
	}
	_, err := base64.StdEncoding.DecodeString(input.PaymentToken)
	if err != nil {
		return nil, exceptions.ErrInvalidPaymentToken
	}
	existingUser, err := os.userGateway.GetUser(input.UserID)
	if err != nil {
		return nil, err
	}
	order := domain.NewOrder()
	for _, inputItem := range input.Items {
		existingItem, err := os.itemRepository.GetItem(inputItem.ID)
		if err != nil {
			return nil, err
		}
		if existingItem == nil {
			return nil, exceptions.ErrItemNotFound
		}
		order.AddItem(existingItem, inputItem.Quantity)
	}
	total := order.CalculateTotalAmount()
	if err := os.orderRepository.SaveOrder(order); err != nil {
		return nil, err
	}
	makedOrderEvent := map[string]any{
		"customer": map[string]any{
			"id":    existingUser["id"],
			"name":  existingUser["name"],
			"email": existingUser["email"],
		},
		"order": map[string]any{
			"order_id": order.ID,
			"total":    total,
		},
		"payment": input.PaymentToken,
	}
	b, err := json.Marshal(&makedOrderEvent)
	if err != nil {
		return nil, err
	}
	queueURL := "https://localhost.localstack.cloud:4566/000000000000/maked-order"
	if err := os.queue.Publish(queueURL, string(b)); err != nil {
		fmt.Println(err)
		return nil, err
	}
	os.logger.Info(
		"make order",
		map[string]any{
			"order_id":    order.ID,
			"customer_id": existingUser["id"],
		},
	)
	return map[string]any{
		"order_id":   order.ID,
		"total":      total,
		"user_name":  existingUser["name"],
		"user_email": existingUser["email"],
	}, nil
}

func (os *OrderService) UpdateStatusOrder(input map[string]any) error {
	orderID := fmt.Sprintf("%s", input["order_id"])
	paymenID := fmt.Sprintf("%s", input["payment_id"])
	paymentStatus := fmt.Sprintf("%s", input["payment_status"])

	existingOrder, err := os.orderRepository.GetOrder(orderID)
	if err != nil {
		return err
	}
	if paymentStatus == "APROVED PAYMENT" {
		existingOrder.ConfirmOrder()
	} else {
		existingOrder.CalcelOrder()
	}
	if err := os.orderRepository.UpdateOrder(existingOrder); err != nil {
		return err
	}
	os.logger.Info(
		"update status order",
		map[string]any{
			"payment_id":     paymenID,
			"order_id":       existingOrder.ID,
			"payment_status": paymentStatus,
			"order_status":   existingOrder.Status,
			"total":          existingOrder.Total,
		},
	)
	return nil
}

func (os *OrderService) GetOrder(orderID string) (map[string]any, error) {
	order, err := os.orderRepository.GetOrder(orderID)
	if err != nil {
		return nil, err
	}
	if order == nil {
		return nil, exceptions.ErrOrderNotFound
	}
	output := map[string]any{
		"order_id":    orderID,
		"order_items": order.OrderItems,
		"total":       order.Total,
		"status":      order.Status,
	}
	return output, nil
}
