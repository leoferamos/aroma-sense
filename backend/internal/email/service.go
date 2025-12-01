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

	// SendAccountDeactivated sends notification when account is deactivated
	SendAccountDeactivated(to, reason string, contestationDeadline string) error

	// SendContestationReceived sends confirmation when contestation is received
	SendContestationReceived(to string) error

	// SendContestationResult sends the result of contestation review
	SendContestationResult(to string, approved bool, reason string) error

	// SendDeletionRequested notifies user that their deletion request was received
	SendDeletionRequested(to string, cancelLink string) error

	// SendDeletionAutoConfirmed notifies user that their deletion request was auto-confirmed
	SendDeletionAutoConfirmed(to string) error

	// SendDeletionCancelled notifies user that their deletion request was cancelled
	SendDeletionCancelled(to string) error

	// SendDataAnonymized notifies user that their personal data has been anonymized
	SendDataAnonymized(to string) error
}
