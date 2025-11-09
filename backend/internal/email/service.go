package email

import "github.com/leoferamos/aroma-sense/internal/model"

// EmailService defines the interface for sending emails.
type EmailService interface {
	// SendPasswordResetCode sends a 6-digit code to reset user's password
	SendPasswordResetCode(to, code string) error

	// SendOrderConfirmation sends order confirmation email to customer
	SendOrderConfirmation(to string, order *model.Order) error

	// SendWelcomeEmail sends welcome email to new users
	SendWelcomeEmail(to, name string) error

	// SendPromotional sends promotional/marketing emails
	SendPromotional(to, subject, htmlBody string) error
}
