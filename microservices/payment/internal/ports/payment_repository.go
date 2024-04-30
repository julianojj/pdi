package ports

import "pdi/payment/internal/core/domain"

type PaymentRepository interface {
	Save(payment *domain.Payment) error
}
