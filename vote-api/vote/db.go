package vote

import (
	"database/sql"
	_ "github.com/lib/pq"
	"log/slog"
	"os"
)

const (
	PgUrlEnvVar          = "postgres://vote_user:password123@database:5432/vote_db"
	CreateVoteTableQuery = `
		CREATE TABLE IF NOT EXISTS votes (
		    vote_id SERIAL PRIMARY KEY,
			imdb_id VARCHAR(255),
			vote_type VARCHAR(255)
		);
	`
	InsertVoteQuery       = "INSERT INTO votes (imdb_id, vote_type) VALUES ($1, $2)"
	GetVoteByMovieIdQuery = "SELECT vote_type, count(*) FROM votes WHERE imdb_id = $1 GROUP BY vote_type"
)

var PgUrl = os.Getenv(PgUrlEnvVar)

func createDbConnection() (*sql.DB, error) {

	if PgUrl == "" {
		panic("PG_URL env var must be set")
	}

	return sql.Open("postgres", PgUrl)
}

func SetupDb() error {
	db, err := createDbConnection()
	if err != nil {
		return err
	}
	defer db.Close()

	_, err = db.Exec(CreateVoteTableQuery)
	return err
}

func InsertVote(vote Vote) error {
	slog.Info("insert vote", slog.String("movie", vote.ImdbId), slog.String("voteType", vote.VoteType))
	db, err := createDbConnection()
	if err != nil {
		return err
	}
	defer db.Close()

	_, err = db.Exec(InsertVoteQuery, vote.ImdbId, vote.VoteType)
	return err
}

func GetVoteByMovieId(movieId string) (map[string]int, error) {
	db, err := createDbConnection()
	if err != nil {
		return nil, err
	}
	defer db.Close()

	rows, err := db.Query(GetVoteByMovieIdQuery, movieId)
	if err != nil {
		return nil, err
	}

	defer rows.Close()

	votes := make(map[string]int)
	for rows.Next() {
		var voteType string
		var voteCount int
		err := rows.Scan(&voteType, &voteCount)
		if err != nil {
			return nil, err
		}
		votes[voteType] = voteCount
	}

	return votes, nil
}
