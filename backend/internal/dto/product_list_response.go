package dto

// ProductListResponse represents a paginated envelope for product search results
type ProductListResponse struct {
	Items []ProductResponse `json:"items"`
	Total int               `json:"total"`
	Page  int               `json:"page"`
	Limit int               `json:"limit"`
}
