package main

import (
	"fmt"
	"github.com/Anthony-Jhoiro/devops-tp-final/vote-api/movies"
	"github.com/Anthony-Jhoiro/devops-tp-final/vote-api/vote"
	"log/slog"
	"net/http"
	"os"
)

const (
	port          = "8080"
	JsonLogEnvVar = "JSON_LOG"
)

func init() {
	if os.Getenv(JsonLogEnvVar) == "true" {
		slog.SetDefault(slog.New(slog.NewJSONHandler(os.Stdout, nil)))
	}
}

func main() {

	if err := vote.SetupDb(); err != nil {
		slog.Error("fail to setup database", slog.String("message", err.Error()))
		return
	}

	http.HandleFunc("/movies", movies.HandleRequests)
	http.HandleFunc("/votes", vote.HandleRequests)

	slog.Info("start application on port", slog.String("port", port))
	err := http.ListenAndServe(fmt.Sprintf(":%s", port), nil)

	if err != nil {
		slog.Error("fail to start application", slog.String("message", err.Error()))
	}
}