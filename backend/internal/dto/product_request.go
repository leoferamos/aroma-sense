package dto

// ProductFormDTO represents the expected payload for creating a product.
type ProductFormDTO struct {
	Name          string   `form:"name" binding:"required"`
	Brand         string   `form:"brand" binding:"required"`
	Weight        float64  `form:"weight" binding:"required"`
	Description   string   `form:"description"`
	Price         float64  `form:"price" binding:"required"`
	Category      string   `form:"category" binding:"required"`
	Notes         []string `form:"notes" binding:"required"`
	StockQuantity int      `form:"stock_quantity" binding:"required,gte=0"`
}

// UpdateProductRequest represents the payload for updating a product.
type UpdateProductRequest struct {
	Name          *string  `json:"name,omitempty"`
	Brand         *string  `json:"brand,omitempty"`
	Weight        *float64 `json:"weight,omitempty"`
	Description   *string  `json:"description,omitempty"`
	Price         *float64 `json:"price,omitempty"`
	Category      *string  `json:"category,omitempty"`
	Notes         []string `json:"notes,omitempty"`
	StockQuantity *int     `json:"stock_quantity,omitempty"`
}
