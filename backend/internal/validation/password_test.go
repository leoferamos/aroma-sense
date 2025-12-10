package validation

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestValidatePassword(t *testing.T) {
	t.Parallel()

	cases := []struct {
		name  string
		pw    string
		email string
		ok    bool
	}{
		{name: "valid basic", pw: "Abcdef1G", email: "user@example.com", ok: true},
		{name: "too short", pw: "Ab1g", email: "user@example.com", ok: false},
		{name: "no upper", pw: "abcdef12", email: "user@example.com", ok: false},
		{name: "no lower", pw: "ABCDEFG1", email: "user@example.com", ok: false},
		{name: "no digit", pw: "Abcdefgh", email: "user@example.com", ok: false},
	}

	for _, tc := range cases {
		tc := tc
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()

			err := ValidatePassword(tc.pw, tc.email)
			if tc.ok {
				require.NoError(t, err)
				return
			}
			require.Error(t, err)
		})
	}
}
