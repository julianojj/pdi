package service

import (
	"time"

	lSdk "github.com/julianojj/essentials-sdk-go/pkg/logger"
	"github.com/julianojj/pdi/notification/internal/ports"
)

type (
	NotificationService struct {
		logger                 lSdk.Logger
		notificationRepository ports.NotificationRepository
	}
)

func NewNotificationService(
	logger lSdk.Logger,
	notificationRepository ports.NotificationRepository,
) *NotificationService {
	return &NotificationService{
		logger,
		notificationRepository,
	}
}

func (ns *NotificationService) NotifyPaymentOrder(input map[string]any) error {
	message := "Send message to"
	ns.logger.Info(message, input)
	notification := map[string]any{
		"message":  message,
		"customer": input["customer"],
		"send_at":  time.Now().UTC(),
	}
	if err := ns.notificationRepository.Save(notification); err != nil {
		return err
	}
	return nil
}
