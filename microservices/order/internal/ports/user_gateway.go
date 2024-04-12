package ports

type UserGateway interface {
	GetUser(userID string) (map[string]any, error)
}
