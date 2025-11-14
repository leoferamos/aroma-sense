package shipping

// quoteRequest is a generic payload understood by the target API.
type quoteRequest struct {
	From struct {
		PostalCode string `json:"postal_code"`
	} `json:"from"`
	To struct {
		PostalCode string `json:"postal_code"`
	} `json:"to"`
	Services string `json:"services"`
	Options  struct {
		OwnHand           bool    `json:"own_hand"`
		Receipt           bool    `json:"receipt"`
		InsuranceValue    float64 `json:"insurance_value"`
		UseInsuranceValue bool    `json:"use_insurance_value"`
	} `json:"options"`
	Package struct {
		Weight float64 `json:"weight"`
		Height float64 `json:"height"`
		Width  float64 `json:"width"`
		Length float64 `json:"length"`
	} `json:"package"`
}

// providerQuote mirrors the provider's quote item.
type providerQuote struct {
	ID           int     `json:"id"`
	Name         string  `json:"name"`
	Price        float64 `json:"price"`
	DeliveryTime int     `json:"delivery_time"`
	HasError     bool    `json:"has_error"`
	Company      struct {
		Name string `json:"name"`
	} `json:"company"`
}
