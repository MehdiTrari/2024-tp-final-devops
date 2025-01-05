package vote

import (
	"os"
	"testing"
)

func InitializeTestDatabase() {
	err := SetupDb()
	if err != nil {
		panic("Failed to set up database for tests: " + err.Error())
	}
}

func TestMain(m *testing.M) {
	// Initialiser la base de données pour les tests
	err := SetupDb()
	if err != nil {
		panic("Failed to set up database for tests: " + err.Error())
	}

	// Exécuter les tests
	os.Exit(m.Run())
}
