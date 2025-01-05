package movies

import (
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestHandleRequestsMovies(t *testing.T) {
	req := httptest.NewRequest(http.MethodGet, "/movies", nil)
	rec := httptest.NewRecorder()

	HandleRequests(rec, req)

	if rec.Code != http.StatusOK {
		t.Errorf("Expected status 200, got %d", rec.Code)
	}

	// Test si les résultats contiennent des films
	if len(rec.Body.String()) == 0 {
		t.Errorf("Expected non-empty response body")
	}
}
