package vote

import (
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
)

func TestHandlePost(t *testing.T) {
	// Préparez une requête POST
	body := `{"imdbId":"tt1234567", "voteType":"like"}`
	req := httptest.NewRequest(http.MethodPost, "/votes", strings.NewReader(body))
	req.Header.Set("Content-Type", "application/json")

	// Mock response recorder
	rec := httptest.NewRecorder()

	// Action
	HandleRequests(rec, req)

	// Assertions
	if rec.Code != http.StatusOK {
		t.Errorf("Expected status 200, got %d", rec.Code)
	}
}

func TestHandleGet(t *testing.T) {
	// Préparez une base de données en mémoire ou mockée
	err := SetupDb()
	if err != nil {
		t.Fatalf("Failed to set up database: %v", err)
	}

	// Insérez un vote pour tester
	err = InsertVote(Vote{ImdbId: "tt1234567", VoteType: "like"})
	if err != nil {
		t.Fatalf("Failed to insert vote: %v", err)
	}

	// Préparez une requête GET
	req := httptest.NewRequest(http.MethodGet, "/votes?imdbId=tt1234567", nil)
	rec := httptest.NewRecorder()

	// Action
	HandleRequests(rec, req)

	// Assertions
	if rec.Code != http.StatusOK {
		t.Errorf("Expected status 200, got %d", rec.Code)
	}

	// Vérifiez le contenu JSON
	expected := `{"like":1}`
	if rec.Body.String() != expected {
		t.Errorf("Expected %s, got %s", expected, rec.Body.String())
	}
}
