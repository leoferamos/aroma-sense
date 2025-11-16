package dto

import "time"

// ReviewRequest represents the payload to create a review
type ReviewRequest struct {
	Rating  int    `json:"rating" binding:"required,min=1,max=5"`
	Comment string `json:"comment" binding:"max=500"`
}

// ReviewResponse represents a published review returned to clients
type ReviewResponse struct {
	ID            string    `json:"id"`
	Rating        int       `json:"rating"`
	Comment       string    `json:"comment"`
	AuthorDisplay string    `json:"author_display"`
	CreatedAt     time.Time `json:"created_at"`
}

// ReviewListResponse is a paginated list of reviews
type ReviewListResponse struct {
	Items []ReviewResponse `json:"items"`
	Total int              `json:"total"`
	Page  int              `json:"page"`
	Limit int              `json:"limit"`
}

// ReviewSummary aggregates ratings for a product
type ReviewSummary struct {
	Average      float64     `json:"average"`
	Count        int         `json:"count"`
	Distribution map[int]int `json:"distribution"`
}
