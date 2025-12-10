package validation

import "testing"

func TestValidatePassword(t *testing.T) {
	cases := []struct {
		name  string
		pw    string
		email string
		ok    bool
	}{
		{"valid basic", "Abcdef1G", "user@example.com", true},
		{"too short", "Ab1g", "user@example.com", false},
		{"no upper", "abcdef12", "user@example.com", false},
		{"no lower", "ABCDEFG1", "user@example.com", false},
		{"no digit", "Abcdefgh", "user@example.com", false},
	}
	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			err := ValidatePassword(tc.pw, tc.email)
			if tc.ok && err != nil {
				t.Fatalf("expected success got error %v", err)
			}
			if !tc.ok && err == nil {
				t.Fatalf("expected failure got success")
			}
		})
	}
}
