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
		err := ValidatePassword(tc.pw, tc.email)
		if tc.ok && err != nil {
			// report unexpected error
			t.Errorf("%s expected success got error %v", tc.name, err)
		}
		if !tc.ok && err == nil {
			t.Errorf("%s expected failure got success", tc.name)
		}
	}
}
