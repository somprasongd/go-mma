package notification_contracts

import "go-mma/shared/common/registry"

const (
	NotificationServiceKey registry.ServiceKey = "NotificationService"
)

type NotificationService interface {
	SendEmail(to string, subject string, payload map[string]any) error
}
