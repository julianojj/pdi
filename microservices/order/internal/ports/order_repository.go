package ports

import "pdi/order/internal/core/domain"

type (
	OrderRepository interface {
		SaveOrder(order *domain.Order) error
		GetOrder(orderID string) (*domain.Order, error)
		UpdateOrder(order *domain.Order) error
	}
)
