package dto

// ProductFormDTO represents the expected payload for creating a product.
type ProductFormDTO struct {
	Name          string   `form:"name" binding:"required"`
	Brand         string   `form:"brand" binding:"required"`
	Weight        float64  `form:"weight" binding:"required"`
	Description   string   `form:"description"`
	Price         float64  `form:"price" binding:"required"`
	Category      string   `form:"category" binding:"required"`
	StockQuantity int      `form:"stock_quantity" binding:"required,gte=0"`
	Accords       []string `form:"accords"`
	Occasions     []string `form:"occasions"`
	Seasons       []string `form:"seasons"`
	Intensity     string   `form:"intensity"`
	Gender        string   `form:"gender"`
	PriceRange    string   `form:"price_range"`
	NotesTop      []string `form:"notes_top"`
	NotesHeart    []string `form:"notes_heart"`
	NotesBase     []string `form:"notes_base"`
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
	StockQuantity *int     `json:"stock_quantity,omitempty" example:"25"`
	Accords       []string `json:"accords,omitempty"`
	Occasions     []string `json:"occasions,omitempty"`
	Seasons       []string `json:"seasons,omitempty"`
	Intensity     *string  `json:"intensity,omitempty"`
	Gender        *string  `json:"gender,omitempty"`
	PriceRange    *string  `json:"price_range,omitempty"`
	NotesTop      []string `json:"notes_top,omitempty"`
	NotesHeart    []string `json:"notes_heart,omitempty"`
	NotesBase     []string `json:"notes_base,omitempty"`
}
