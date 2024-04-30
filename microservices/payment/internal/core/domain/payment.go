package domain

import "github.com/google/uuid"

type Payment struct {
	ID         string
	OrderID    string
	CustomerID string
	Type       string
	Amount     float64
	Status     string
}

func NewPayment(
	orderID string,
	customerID string,
	paymentType string,
	amount float64,
) *Payment {
	return &Payment{
		ID:         uuid.NewString(),
		OrderID:    orderID,
		CustomerID: customerID,
		Type:       paymentType,
		Amount:     amount,
		Status:     "PENDING_PAYMENT",
	}
}

func (p *Payment) AprovePayment() {
	p.Status = "APROVED PAYMENT"
}

func (p *Payment) RejectPayment() {
	p.Status = "REJECTED PAYMENT"
}
