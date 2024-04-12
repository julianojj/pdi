package domain

import "github.com/google/uuid"

type (
	Order struct {
		ID         string
		OrderItems []*OrderItem
		Total      float64
		Status     string
	}
	OrderItem struct {
		ID       string
		Amount   float64
		Quantity int
	}
)

var (
	PENDING_ORDER_STATUS   = "PENDING_PAYMENT"
	CONFIRMED_ORDER_STATUS = "CONFIRMED_ORDER"
	REJECTED_ORDER_STATUS  = "REJECTED_ORDER"
)

func NewOrder() *Order {
	return &Order{
		ID:         uuid.NewString(),
		OrderItems: make([]*OrderItem, 0),
		Status:     PENDING_ORDER_STATUS,
		Total:      0,
	}
}

func (o *Order) AddItem(item *Item, quantity int) {
	o.OrderItems = append(o.OrderItems, &OrderItem{
		ID:       item.ID,
		Amount:   item.Amount,
		Quantity: quantity,
	})
}

func (o *Order) CalculateTotalAmount() float64 {
	var total float64 = 0
	for _, orderItem := range o.OrderItems {
		total += orderItem.Amount * float64(orderItem.Quantity)
	}
	o.Total = total
	return total
}

func (o *Order) ConfirmOrder() {
	o.Status = CONFIRMED_ORDER_STATUS
}

func (o *Order) CalcelOrder() {
	o.Status = REJECTED_ORDER_STATUS
}
