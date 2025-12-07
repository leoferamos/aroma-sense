package dto

import "time"

// AdminOrderItem represents the data for an order shown in admin listings
type AdminOrderItem struct {
	ID          uint      `json:"id" example:"1"`
	PublicID    string    `json:"public_id" example:"uuid"`
	UserID      string    `json:"user_id" example:"uuid"`
	TotalAmount float64   `json:"total_amount" example:"123.45"`
	Status      string    `json:"status" example:"pending"`
	CreatedAt   time.Time `json:"created_at"`
}

// PaginationMeta contains pagination information
type PaginationMeta struct {
	Page       int `json:"page" example:"1"`
	PerPage    int `json:"per_page" example:"25"`
	TotalPages int `json:"total_pages" example:"10"`
	TotalCount int `json:"total_count" example:"247"`
}

// StatsMeta contains aggregated statistics for the listing
type StatsMeta struct {
	TotalRevenue      float64 `json:"total_revenue" example:"12345.67"`
	AverageOrderValue float64 `json:"average_order_value" example:"49.95"`
}

// AdminOrdersResponse is the response returned by GET /admin/orders
type AdminOrdersResponse struct {
	Orders []AdminOrderItem `json:"orders"`
	Meta   struct {
		Pagination PaginationMeta `json:"pagination"`
		Stats      StatsMeta      `json:"stats"`
	} `json:"meta"`
}
