package service

import "log/slog"

type (
	NotificationService struct {
		logger *slog.Logger
	}
)

func NewNotificationService(
	logger *slog.Logger,
) *NotificationService {
	return &NotificationService{
		logger,
	}
}

func (ns *NotificationService) NotifyPaymentOrder(input map[string]any) error {
	ns.logger.Info(
		"Send message to",
		"data", input,
	)
	return nil
}
