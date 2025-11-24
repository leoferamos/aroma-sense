package dto

// RecommendRequest is the input payload for the AI recommend endpoint.
type RecommendRequest struct {
	Message string   `json:"message"`
	History []string `json:"history,omitempty"`
	Limit   int      `json:"limit,omitempty"`
}

// RecommendSuggestion is a compact product card to show inside the chat.
type RecommendSuggestion struct {
	ID           uint    `json:"id"`
	Name         string  `json:"name"`
	Brand        string  `json:"brand"`
	Slug         string  `json:"slug"`
	ThumbnailURL string  `json:"thumbnail_url"`
	Price        float64 `json:"price"`
	Reason       string  `json:"reason"`
}

// RecommendResponse bundles suggestions and lightweight reasoning.
type RecommendResponse struct {
	Suggestions []RecommendSuggestion `json:"suggestions"`
	Reasoning   string                `json:"reasoning,omitempty"`
}
