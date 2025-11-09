package email

import (
	"fmt"
	"net/smtp"

	"github.com/leoferamos/aroma-sense/internal/model"
)

// SMTPEmailService implements EmailService using SMTP protocol
type SMTPEmailService struct {
	config *SMTPConfig
}

// NewSMTPEmailService creates a new SMTP email service
func NewSMTPEmailService(config *SMTPConfig) (*SMTPEmailService, error) {
	if err := config.Validate(); err != nil {
		return nil, fmt.Errorf("invalid SMTP config: %w", err)
	}
	return &SMTPEmailService{config: config}, nil
}

// SendPasswordResetCode sends a password reset code via email
func (s *SMTPEmailService) SendPasswordResetCode(to, code string) error {
	subject := "Password Reset Code - Aroma Sense"
	htmlBody := PasswordResetTemplate(code)

	return s.sendEmail(to, subject, htmlBody)
}

// SendOrderConfirmation sends order confirmation email
func (s *SMTPEmailService) SendOrderConfirmation(to string, order *model.Order) error {
	subject := "Order Confirmation - Aroma Sense"
	htmlBody := OrderConfirmationTemplate(fmt.Sprintf("#%d", order.ID))

	return s.sendEmail(to, subject, htmlBody)
}

// SendWelcomeEmail sends welcome email to new users
func (s *SMTPEmailService) SendWelcomeEmail(to, name string) error {
	subject := "Welcome to Aroma Sense!"
	htmlBody := WelcomeEmailTemplate(name)

	return s.sendEmail(to, subject, htmlBody)
}

// SendPromotional sends promotional emails
func (s *SMTPEmailService) SendPromotional(to, subject, htmlBody string) error {
	return s.sendEmail(to, subject, htmlBody)
}

// sendEmail is a helper to send an email via SMTP
func (s *SMTPEmailService) sendEmail(to, subject, htmlBody string) error {
	// Build email message with headers
	msg := []byte(fmt.Sprintf(
		"From: %s\r\n"+
			"To: %s\r\n"+
			"Subject: %s\r\n"+
			"MIME-Version: 1.0\r\n"+
			"Content-Type: text/html; charset=UTF-8\r\n"+
			"\r\n"+
			"%s",
		s.config.From,
		to,
		subject,
		htmlBody,
	))

	// Send email via SMTP
	err := smtp.SendMail(
		s.config.Address(),
		s.config.Auth(),
		s.config.From,
		[]string{to},
		msg,
	)

	if err != nil {
		return fmt.Errorf("failed to send email to %s: %w", to, err)
	}

	return nil
}
