package dto

// ShippingOption represents a shipping option returned to the client.
type ShippingOption struct {
	Carrier       string  `json:"carrier"`
	ServiceCode   string  `json:"service_code"`
	Price         float64 `json:"price"`
	EstimatedDays int     `json:"estimated_days"`
}
