package adapters

import (
	"math/rand"
	"pdi/payment/internal/ports"
)

type Pagarme struct{}

func NewPagarme() ports.PaymentGateway {
	return &Pagarme{}
}

func (p *Pagarme) ProcessPayment(map[string]any) (map[string]any, error) {
	var codes []string = []string{"0000", "111"}
	randIndex := rand.Intn(2)
	return map[string]any{
		"code":   codes[randIndex],
		"result": "success to process payment",
	}, nil
}
