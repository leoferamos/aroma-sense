package notification

import (
	"github.com/leoferamos/aroma-sense/internal/email"
	"github.com/leoferamos/aroma-sense/internal/model"
)

// NotificationService defines high-level notifications used by services.
type NotificationService interface {
	SendPasswordResetCode(to, code string) error
	SendWelcomeEmail(to, name string) error
	SendOrderConfirmation(to string, order *model.Order) error
	SendAccountDeactivated(to, reason string, contestationDeadline string) error
	SendContestationReceived(to string) error
	SendContestationResult(to string, approved bool, reason string) error
	SendDeletionRequested(to string, cancelLink string) error
	SendDeletionAutoConfirmed(to string) error
	SendDataAnonymized(to string) error
	SendPromotional(to, subject, htmlBody string) error
}

type notifier struct {
	es           email.EmailService
	frontendBase string
}

// NewNotifier creates a notification service that delegates to the provided EmailService
// frontendBase is optional and used to build frontend links (e.g. cancel link) when callers
// pass an empty link.
func NewNotifier(es email.EmailService, frontendBase string) NotificationService {
	return &notifier{es: es, frontendBase: frontendBase}
}

func (n *notifier) SendPasswordResetCode(to, code string) error {
	return n.es.SendPasswordResetCode(to, code)
}

func (n *notifier) SendWelcomeEmail(to, name string) error {
	return n.es.SendWelcomeEmail(to, name)
}

func (n *notifier) SendOrderConfirmation(to string, order *model.Order) error {
	return n.es.SendOrderConfirmation(to, order)
}

func (n *notifier) SendAccountDeactivated(to, reason string, contestationDeadline string) error {
	return n.es.SendAccountDeactivated(to, reason, contestationDeadline)
}

func (n *notifier) SendContestationReceived(to string) error {
	return n.es.SendContestationReceived(to)
}

func (n *notifier) SendContestationResult(to string, approved bool, reason string) error {
	return n.es.SendContestationResult(to, approved, reason)
}

func (n *notifier) SendDeletionRequested(to string, cancelLink string) error {
	// Build a cancel link if blank
	if cancelLink == "" {
		if n.frontendBase != "" {
			cancelLink = n.frontendBase + "/settings/account"
		} else {
			cancelLink = "/account/settings"
		}
	}
	return n.es.SendDeletionRequested(to, cancelLink)
}

func (n *notifier) SendDeletionAutoConfirmed(to string) error {
	return n.es.SendDeletionAutoConfirmed(to)
}

func (n *notifier) SendDataAnonymized(to string) error {
	return n.es.SendDataAnonymized(to)
}

func (n *notifier) SendPromotional(to, subject, htmlBody string) error {
	return n.es.SendPromotional(to, subject, htmlBody)
}
