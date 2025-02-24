package features

import (
	"context"
	"go-mma/shared/common/mediator"
	"log"
)

type SendEmailCommand struct {
	To      string                 `json:"to"`
	Subject string                 `json:"subject"`
	Payload map[string]interface{} `json:"payload"`
}

type sendEmailHandler struct {
}

func NewSendEmailHandler() *sendEmailHandler {
	return &sendEmailHandler{}
}

func (h *sendEmailHandler) Handle(ctx context.Context, cmd *SendEmailCommand) (*mediator.NoResponse, error) {
	// Implement email sending logic here
	log.Println("Sending email to", cmd.To, "with subject:", cmd.Subject, "and payload:", cmd.Payload)
	return &mediator.NoResponse{}, nil
}
