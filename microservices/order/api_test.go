package main

import (
	"encoding/base64"
	"encoding/json"
	"fmt"
	"testing"
	"time"

	"github.com/go-resty/resty/v2"
	"github.com/stretchr/testify/assert"
)

func TestAPI(t *testing.T) {
	payment := map[string]any{
		"method":        "credit",
		"brand":         "mastercard",
		"card_number":   "5188129146451687",
		"security_code": "851",
		"validate":      "2025-10-20",
		"installments":  "1",
	}
	b, _ := json.Marshal(&payment)
	paymentToken := base64.StdEncoding.EncodeToString(b)
	fmt.Println(paymentToken)
	input := map[string]any{
		"user_id":       "f5a0a785-75d2-4d8b-ae33-0039ee24216f",
		"payment_token": paymentToken,
		"items": []map[string]any{
			{
				"item_id":  "1",
				"quantity": 2,
			},
			{
				"item_id":  "2",
				"quantity": 2,
			},
		},
	}
	output := map[string]any{}
	resty.New().R().SetBody(input).SetResult(&output).Post("http://localhost:8081/orders")
	time.Sleep(time.Second)
	outputOrder := map[string]any{}
	resty.New().R().SetResult(&outputOrder).Get(fmt.Sprintf("http://localhost:8081/orders/%s", output["order_id"]))
	assert.NotNil(t, output["order_id"])
	assert.Equal(t, 2000.00, output["total"])
	assert.Equal(t, "Juliano Silva", output["user_name"])
	assert.Equal(t, "juliano.silva@stone.com.br", output["user_email"])
	assert.Equal(t, "CONFIRMED_PAYMENT", outputOrder["status"])
}
