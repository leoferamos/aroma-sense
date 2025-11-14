package model

// Parcel represents a package to be shipped.
type Parcel struct {
	WeightKg float64 `json:"weight_kg"`
	LengthCm float64 `json:"length_cm"`
	WidthCm  float64 `json:"width_cm"`
	HeightCm float64 `json:"height_cm"`
}

// Shipment captures shipping metadata in the domain layer.
type Shipment struct {
	Carrier       string  `json:"carrier"`
	ServiceCode   string  `json:"service_code"`
	Price         float64 `json:"price"`
	EstimatedDays int     `json:"estimated_days"`
	Tracking      string  `json:"tracking,omitempty"`
	Status        string  `json:"status,omitempty"`
}
