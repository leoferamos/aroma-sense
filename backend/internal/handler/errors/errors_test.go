package handlererrors

import (
	"errors"
	"net/http"
	"testing"

	"github.com/leoferamos/aroma-sense/internal/apperror"
	"github.com/stretchr/testify/require"
)

func TestMapServiceError(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name   string
		input  error
		ok     bool
		status int
		code   string
	}{
		{name: "domain code mapped", input: apperror.NewCodeMessage("unauthenticated", "unauthenticated"), ok: true, status: http.StatusUnauthorized, code: "unauthenticated"},
		{name: "domain unknown code defaults 500", input: apperror.NewCodeMessage("unknown_code", "oops"), ok: true, status: http.StatusInternalServerError, code: "unknown_code"},
		{name: "domain empty code becomes internal", input: apperror.NewDomain(errors.New("boom"), "", ""), ok: true, status: http.StatusInternalServerError, code: "internal_error"},
		{name: "nil error", input: nil, ok: false},
		{name: "non domain error", input: errors.New("plain"), ok: false},
	}

	for _, tt := range tests {
		caseData := tt
		t.Run(caseData.name, func(t *testing.T) {
			t.Parallel()

			status, msg, ok := mapServiceError(caseData.input)
			require.Equal(t, caseData.ok, ok)

			if !caseData.ok {
				return
			}
			require.Equal(t, caseData.status, status)
			require.Equal(t, caseData.code, msg)
		})
	}
}
