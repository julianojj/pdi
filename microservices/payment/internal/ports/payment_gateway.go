package ports

type PaymentGateway interface {
	ProcessPayment(map[string]any) (map[string]any, error)
}