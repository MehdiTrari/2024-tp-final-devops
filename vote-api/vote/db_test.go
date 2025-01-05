package vote

import (
	"testing"
)

func TestInsertVote(t *testing.T) {
	// Préparez une base de données en mémoire ou mockée (utilisez une vraie DB pour des tests locaux si nécessaire)
	err := SetupDb()
	if err != nil {
		t.Fatalf("Failed to set up database: %v", err)
	}

	// Données de test
	vote := Vote{
		ImdbId:   "tt1234567",
		VoteType: "like",
	}

	// Action
	err = InsertVote(vote)

	// Assertions
	if err != nil {
		t.Errorf("InsertVote() failed: %v", err)
	}
}

func TestGetVoteByMovieId(t *testing.T) {
	// Préparez une base de données en mémoire ou mockée
	err := SetupDb()
	if err != nil {
		t.Fatalf("Failed to set up database: %v", err)
	}

	// Insérez un vote
	err = InsertVote(Vote{ImdbId: "tt1234567", VoteType: "like"})
	if err != nil {
		t.Fatalf("InsertVote() failed: %v", err)
	}

	// Action
	votes, err := GetVoteByMovieId("tt1234567")

	// Assertions
	if err != nil {
		t.Errorf("GetVoteByMovieId() failed: %v", err)
	}
	if votes["like"] != 1 {
		t.Errorf("Expected 1 like vote, got %d", votes["like"])
	}
}
