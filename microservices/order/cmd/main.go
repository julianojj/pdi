package main

import (
	"log/slog"
	"net/http"
	"os"
	"pdi/order/internal/adapters"
	"pdi/order/internal/core/service"

	"github.com/gin-gonic/gin"
)

func main() {
	orderRepository := adapters.NewOrderRepositoryDynamoDB()
	itemRepository := adapters.NewItemRepositoryMemory()
	userGateway := adapters.NewUserGatewayAPI()
	queue := adapters.NewSQS()

	logger := slog.New(slog.NewJSONHandler(os.Stdout, nil)).With(
		slog.Any("application", map[string]any{
			"name":        "order-ms",
			"environment": "dev",
			"version":     "1.0.0",
		}),
	)

	orderService := service.NewOrderService(orderRepository, itemRepository, queue, userGateway, logger)

	r := gin.Default()

	r.POST("/orders", func(ctx *gin.Context) {
		var input *service.OrderServiceInput
		if err := ctx.ShouldBindJSON(&input); err != nil {
			ctx.JSON(http.StatusBadRequest, map[string]any{
				"message": err.Error(),
			})
			return
		}
		output, err := orderService.MakeOrder(input)
		if err != nil {
			logger.Error(
				"error to make order",
				slog.Any("data", map[string]any{
					"err": err,
				}),
			)
			ctx.JSON(http.StatusUnprocessableEntity, map[string]any{
				"message": err.Error(),
			})
			return
		}
		ctx.JSON(http.StatusCreated, output)
	})

	r.GET("/orders/:id", func(ctx *gin.Context) {
		output, err := orderService.GetOrder(ctx.Param("id"))
		if err != nil {
			ctx.AbortWithStatusJSON(http.StatusFailedDependency, map[string]any{
				"message": err.Error(),
			})
			return
		}
		ctx.JSON(http.StatusOK, output)
	})

	r.Run(":8081")
}
