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
	StockQuantity int      `form:"stock_quantity" binding:"required"`
}
