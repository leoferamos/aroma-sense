package dto

import "strings"

// CreatePaymentIntentRequest represents payload to start payment.
type CreatePaymentIntentRequest struct {
	ShippingAddress   string             `json:"shipping_address" binding:"required"`
	ShippingSelection *ShippingSelection `json:"shipping_selection,omitempty"`
	CustomerEmail     string             `json:"customer_email,omitempty"`
	OrderPublicID     string             `json:"order_public_id,omitempty"`
}

// ShippingPostalCode attempts to extract the CEP digits from the shipping address.
func (r *CreatePaymentIntentRequest) ShippingPostalCode() string {
	if r == nil {
		return ""
	}
	// naive extract: keep digits only
	clean := ""
	for _, ch := range r.ShippingAddress {
		if ch >= '0' && ch <= '9' {
			clean += string(ch)
		}
	}
	if len(clean) >= 8 {
		return clean[len(clean)-8:]
	}
	return strings.TrimSpace(clean)
}
