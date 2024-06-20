package ports

type NotificationRepository interface {
	Save(notification map[string]any) error
}
