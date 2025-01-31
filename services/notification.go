package services

import "log"

type NotificationService struct {
}

func NewNotificationService() *NotificationService {
	return &NotificationService{}
}

func (s *NotificationService) SendEmail(to string, subject string, payload map[string]any) error {
	// implement email sending logic here
	log.Println("Sending email to", to, "with subject:", subject, "and payload:", payload)
	return nil
}
