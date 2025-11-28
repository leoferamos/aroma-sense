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

	htmlBody := fmt.Sprintf(`
		<h2>Conta Desativada</h2>
		<p>Olá,</p>
		<p>Informamos que sua conta no Aroma Sense foi desativada pelos seguintes motivos:</p>
		<p><strong>%s</strong></p>
		<p>Você tem até <strong>%s</strong> para apresentar contestação através do nosso suporte.</p>
		<p>Para contestar, acesse sua conta ou entre em contato conosco.</p>
		<p>Atenciosamente,<br>Equipe Aroma Sense</p>
	`, reason, contestationDeadline)

	return s.sendEmail(to, subject, htmlBody)
}

// SendContestationReceived sends confirmation when contestation is received
func (s *SMTPEmailService) SendContestationReceived(to string) error {
	subject := "Contestação Recebida - Aroma Sense"

	htmlBody := `
		<h2>Contestação Recebida</h2>
		<p>Olá,</p>
		<p>Recebemos sua contestação sobre a desativação da conta.</p>
		<p>Nossa equipe irá analisar o caso em até 5 dias úteis e entraremos em contato.</p>
		<p>Atenciosamente,<br>Equipe Aroma Sense</p>
	`

	return s.sendEmail(to, subject, htmlBody)
}

// SendContestationResult sends the result of contestation review
func (s *SMTPEmailService) SendContestationResult(to string, approved bool, reason string) error {
	var subject, status string

	if approved {
		subject = "Contestação Aprovada - Conta Reativada"
		status = "aprovada"
	} else {
		subject = "Contestação Rejeitada - Aroma Sense"
		status = "rejeitada"
	}

	htmlBody := fmt.Sprintf(`
		<h2>Resultado da Contestação</h2>
		<p>Olá,</p>
		<p>Sua contestação foi <strong>%s</strong>.</p>
		<p><strong>Motivo:</strong> %s</p>
		<p>Atenciosamente,<br>Equipe Aroma Sense</p>
	`, status, reason)

	return s.sendEmail(to, subject, htmlBody)
}
