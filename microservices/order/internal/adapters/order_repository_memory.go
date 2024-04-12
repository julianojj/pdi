package adapters

import (
	"pdi/order/internal/core/domain"
	"pdi/order/internal/ports"
)

type OrderDataMemory struct {
	orders []*domain.Order
}

func NewOrderDataMemory() ports.OrderRepository {
	orders := make([]*domain.Order, 0)
	return &OrderDataMemory{
		orders,
	}
}

func (o *OrderDataMemory) SaveOrder(order *domain.Order) error {
	o.orders = append(o.orders, order)
	return nil
}

func (o *OrderDataMemory) GetOrder(orderID string) (*domain.Order, error) {
	for _, order := range o.orders {
		if order.ID == orderID {
			return order, nil
		}
	}
	return nil, nil
}

func (o *OrderDataMemory) UpdateOrder(order *domain.Order) error {
	for _, o := range o.orders {
		if o.ID == order.ID {
			o = order
			return nil
		}
	}
	return nil
}
