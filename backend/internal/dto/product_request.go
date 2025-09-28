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
// @Description Product update request (all fields are optional)
type UpdateProductRequest struct {
	Name          *string  `json:"name,omitempty" example:"Sauvage Elixir"`
	Brand         *string  `json:"brand,omitempty" example:"Dior"`
	Weight        *float64 `json:"weight,omitempty" example:"60.0"`
	Description   *string  `json:"description,omitempty" example:"An intense and spicy fragrance"`
	Price         *float64 `json:"price,omitempty" example:"399.99"`
	Category      *string  `json:"category,omitempty" example:"Eau de Parfum"`
	Notes         []string `json:"notes,omitempty" example:"cinnamon,nutmeg,cardamom"`
	StockQuantity *int     `json:"stock_quantity,omitempty" example:"25"`
}
