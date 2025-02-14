package service

import (
	"log"

	notificationsContracts "go-mma/shared/contracts/notification_contracts"
)

type notificationService struct {
}

func NewNotificationService() notificationsContracts.NotificationService {
	return &notificationService{}
}

func (s *notificationService) SendEmail(to string, subject string, payload map[string]any) error {
	// implement email sending logic here
	log.Println("Sending email to", to, "with subject:", subject, "and payload:", payload)
	return nil
}
