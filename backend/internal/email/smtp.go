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

// SendAccountDeactivated sends notification when account is deactivated
func (s *SMTPEmailService) SendAccountDeactivated(to, reason string, contestationDeadline string) error {
	subject := "Sua conta no Aroma Sense foi desativada"
	htmlBody := AccountDeactivatedTemplate(reason, contestationDeadline)
	return s.sendEmail(to, subject, htmlBody)
}

// SendContestationReceived sends confirmation when contestation is received
func (s *SMTPEmailService) SendContestationReceived(to string) error {
	subject := "Contestação Recebida - Aroma Sense"
	htmlBody := ContestationReceivedTemplate()
	return s.sendEmail(to, subject, htmlBody)
}

// SendContestationResult sends the result of contestation review
func (s *SMTPEmailService) SendContestationResult(to string, approved bool, reason string) error {
	subject := "Contestação Resultado - Aroma Sense"
	htmlBody := ContestationResultTemplate(approved, reason)
	return s.sendEmail(to, subject, htmlBody)
}

// SendDeletionRequested notifies the user that their deletion request was received
func (s *SMTPEmailService) SendDeletionRequested(to string, cancelLink string) error {
	subject := "Pedido de exclusão recebido — Aroma Sense"
	// requestedAt not available here; use a generic message
	htmlBody := DeletionRequestedTemplate("", "agora", cancelLink)
	return s.sendEmail(to, subject, htmlBody)
}

// SendDeletionAutoConfirmed notifies the user that their deletion was auto-confirmed
func (s *SMTPEmailService) SendDeletionAutoConfirmed(to string) error {
	subject := "Exclusão da conta confirmada — Aroma Sense"
	htmlBody := DeletionAutoConfirmedTemplate("", "agora")
	return s.sendEmail(to, subject, htmlBody)
}

// SendDataAnonymized notifies the user that their data was anonymized
func (s *SMTPEmailService) SendDataAnonymized(to string) error {
	subject := "Seus dados foram anonimizados — Aroma Sense"
	htmlBody := DataAnonymizedTemplate("agora")
	return s.sendEmail(to, subject, htmlBody)
}

// SendDeletionCancelled notifies the user that their deletion request was cancelled
func (s *SMTPEmailService) SendDeletionCancelled(to string) error {
	subject := "Solicitação de exclusão cancelada — Aroma Sense"
	htmlBody := DeletionCancelledTemplate("", "agora")
	return s.sendEmail(to, subject, htmlBody)
}
