package httpapi

import (
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestHealth(t *testing.T) {
	req := httptest.NewRequest(http.MethodGet, "/health", nil)
	res := httptest.NewRecorder()

	NewHandler("auto-workflow").ServeHTTP(res, req)

	if res.Code != http.StatusOK {
		t.Fatalf("status = %d, want %d", res.Code, http.StatusOK)
	}
	if got := res.Header().Get("Content-Type"); got != "application/json" {
		t.Fatalf("content type = %q, want application/json", got)
	}
	if got := res.Body.String(); got != "{\"service\":\"auto-workflow\",\"status\":\"ok\"}\n" {
		t.Fatalf("body = %q", got)
	}
}

func TestUnknownPath(t *testing.T) {
	req := httptest.NewRequest(http.MethodGet, "/missing", nil)
	res := httptest.NewRecorder()

	NewHandler("auto-workflow").ServeHTTP(res, req)

	if res.Code != http.StatusNotFound {
		t.Fatalf("status = %d, want %d", res.Code, http.StatusNotFound)
	}
}
