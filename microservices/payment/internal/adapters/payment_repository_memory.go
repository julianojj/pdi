package adapters

import (
	"pdi/payment/internal/core/domain"
	"pdi/payment/internal/ports"
)

type PaymentRepositoryMemory struct {
	Payments []*domain.Payment
}

func NewPaymentRepositoryMemory() ports.PaymentRepository {
	return &PaymentRepositoryMemory{
		Payments: make([]*domain.Payment, 0),
	}
}

func (p *PaymentRepositoryMemory) Save(payment *domain.Payment) error {
	p.Payments = append(p.Payments, payment)
	return nil
}
