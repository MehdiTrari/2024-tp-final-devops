package vote

import (
	"os"
	"testing"
)

func TestSetupDb(t *testing.T) {
	err := SetupDb()
	if err != nil {
		t.Fatalf("SetupDb failed: %v", err)
	}
}

func TestMain(m *testing.M) {
	err := SetupDb()
	if err != nil {
		panic("Failed to set up database for tests: " + err.Error())
	}

	// Exécuter les tests
	os.Exit(m.Run())
}
