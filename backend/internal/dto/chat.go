package dto

// ChatRequest represents a conversational turn from the user.
type ChatRequest struct {
	Message   string   `json:"message"`
	SessionID string   `json:"session_id,omitempty"`
	History   []string `json:"history,omitempty"`
}

// ChatResponse is the model reply plus optional product suggestions.
type ChatResponse struct {
	Reply        string                `json:"reply"`
	Suggestions  []RecommendSuggestion `json:"suggestions,omitempty"`
	FollowUpHint string                `json:"follow_up_hint,omitempty"`
}
