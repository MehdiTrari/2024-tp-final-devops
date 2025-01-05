package vote

import (
	"testing"
)

func cleanDb() {
	db, err := createDbConnection()
	if err != nil {
		panic("Failed to connect to the database for cleanup: " + err.Error())
	}
	defer db.Close()

	_, err = db.Exec("DELETE FROM votes")
	if err != nil {
		panic("Failed to clean up the database: " + err.Error())
	}
}

func TestInsertVote(t *testing.T) {
	cleanDb() // Nettoyer la base avant le test

	// Données de test
	vote := Vote{
		ImdbId:   "tt1234567",
		VoteType: "like",
	}

	// Action
	err := InsertVote(vote)

	// Assertions
	if err != nil {
		t.Errorf("InsertVote() failed: %v", err)
	}
}

func TestGetVoteByMovieId(t *testing.T) {
	cleanDb() // Nettoyer la base avant le test

	// Insérez un vote
	err := InsertVote(Vote{ImdbId: "tt1234567", VoteType: "like"})
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
