package handler

import (
	"errors"
	"net/http"
	"testing"

	"github.com/leoferamos/aroma-sense/internal/apperror"
)

func TestMapServiceError_DomainCodeMapped(t *testing.T) {
	de := apperror.NewCodeMessage("unauthenticated", "unauthenticated")
	status, msg, ok := mapServiceError(de)
	if !ok {
		t.Fatalf("expected ok=true")
	}
	if status != http.StatusUnauthorized {
		t.Fatalf("unexpected status: %d", status)
	}
	if msg != "unauthenticated" {
		t.Fatalf("unexpected message: %s", msg)
	}
}

func TestMapServiceError_DomainCodeUnknownKeepsCode500(t *testing.T) {
	de := apperror.NewCodeMessage("unknown_code", "oops")
	status, msg, ok := mapServiceError(de)
	if !ok {
		t.Fatalf("expected ok=true")
	}
	if status != http.StatusInternalServerError {
		t.Fatalf("unexpected status: %d", status)
	}
	if msg != "unknown_code" {
		t.Fatalf("unexpected message: %s", msg)
	}
}

func TestMapServiceError_DomainEmptyCodeDefaultsInternal(t *testing.T) {
	de := apperror.NewDomain(errors.New("boom"), "", "")
	status, msg, ok := mapServiceError(de)
	if !ok {
		t.Fatalf("expected ok=true")
	}
	if status != http.StatusInternalServerError {
		t.Fatalf("unexpected status: %d", status)
	}
	if msg != "internal_error" {
		t.Fatalf("unexpected message: %s", msg)
	}
}

func TestMapServiceError_NilOrNonDomain(t *testing.T) {
	if _, _, ok := mapServiceError(nil); ok {
		t.Fatalf("expected ok=false for nil")
	}
	if _, _, ok := mapServiceError(errors.New("plain")); ok {
		t.Fatalf("expected ok=false for non-domain error")
	}
}
