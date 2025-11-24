package service

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestExtractImageNameFromURL(t *testing.T) {
	tests := []struct {
		name     string
		url      string
		expected string
	}{
		{
			name:     "Valid Supabase URL",
			url:      "https://domain.com/bucket/image-name.jpg",
			expected: "image-name.jpg",
		},
		{
			name:     "Full Supabase storage URL",
			url:      "https://opgdiulsnfcxjymmelny.supabase.co/storage/v1/object/public/Photos/product-54f615e1-fe25-4620-a2f2-a73058c6f61a.jpg",
			expected: "product-54f615e1-fe25-4620-a2f2-a73058c6f61a.jpg",
		},
		{
			name:     "Supabase thumbnail URL",
			url:      "https://opgdiulsnfcxjymmelny.supabase.co/storage/v1/object/public/Photos/product-54f615e1-fe25-4620-a2f2-a73058c6f61a_thumb.jpg",
			expected: "product-54f615e1-fe25-4620-a2f2-a73058c6f61a_thumb.jpg",
		},
		{
			name:     "URL with multiple slashes",
			url:      "https://domain.com/bucket/folder/image-name.png",
			expected: "image-name.png",
		},
		{
			name:     "Empty URL",
			url:      "",
			expected: "",
		},
		{
			name:     "URL without slashes",
			url:      "invalid-url",
			expected: "",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := extractImageNameFromURL(tt.url)
			assert.Equal(t, tt.expected, result)
		})
	}
}
