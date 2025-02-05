package services

import "log"

type NotificationService interface {
	SendEmail(to string, subject string, payload map[string]any) error
}

type notificationService struct {
}

func NewNotificationService() NotificationService {
	return &notificationService{}
}

func (s *notificationService) SendEmail(to string, subject string, payload map[string]any) error {
	// implement email sending logic here
	log.Println("Sending email to", to, "with subject:", subject, "and payload:", payload)
	return nil
}
