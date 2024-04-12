package service

import (
	"encoding/json"
	"errors"
	"fmt"
	"log/slog"
	"pdi/order/internal/core/domain"
	"pdi/order/internal/ports"
)

type (
	OrderService struct {
		orderRepository ports.OrderRepository
		itemRepository  ports.ItemRepository
		queue           ports.Queue
		userGateway     ports.UserGateway
		logger          *slog.Logger
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
	queue ports.Queue,
	userGateway ports.UserGateway,
	logger *slog.Logger,
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
		return nil, errors.New("user id is required")
	}
	if input.PaymentToken == "" {
		return nil, errors.New("payment token is required")
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
			return nil, errors.New("item not found")
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
	if err := os.queue.Publish(string(b)); err != nil {
		fmt.Println(err)
		return nil, err
	}
	os.logger.Info(
		"make order",
		slog.Any("data", map[string]any{
			"order_id":    order.ID,
			"customer_id": existingUser["id"],
		}),
		slog.String("path", "service.order_service.go"),
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
		slog.Any("data", map[string]any{
			"payment_id":     paymenID,
			"order_id":       existingOrder.ID,
			"payment_status": paymentStatus,
			"order_status":   existingOrder.Status,
			"total":          existingOrder.Total,
		}),
		slog.String("path", "service.order_service.UpdateStatusOrder"),
	)
	return nil
}

func (os *OrderService) GetOrder(orderID string) (map[string]any, error) {
	order, err := os.orderRepository.GetOrder(orderID)
	if err != nil {
		return nil, err
	}
	if order == nil {
		return nil, errors.New("order not found")
	}
	output := map[string]any{
		"order_id":    orderID,
		"order_items": order.OrderItems,
		"total":       order.Total,
		"status":      order.Status,
	}
	return output, nil
}
