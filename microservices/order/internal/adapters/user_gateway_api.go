package adapters

import (
	"errors"
	"fmt"
	"pdi/order/internal/ports"

	"github.com/go-resty/resty/v2"
)

type UserGatewayAPI struct {
	httpClient *resty.Client
}

func NewUserGatewayAPI() ports.UserGateway {
	return &UserGatewayAPI{
		httpClient: resty.New(),
	}
}

func (u *UserGatewayAPI) GetUser(userID string) (map[string]any, error) {
	var url = fmt.Sprintf("http://user_api:8080/users/%s", userID)
	var output map[string]any
	response, err := u.httpClient.
		R().
		SetResult(&output).
		Get(url)
	if err != nil {
		return output, err
	}
	if response.StatusCode() != 200 {
		return output, errors.New("error to get user")
	}
	return output, nil
}
