package router

import (
	"net/http"

	"dotaWorstPlayerChacker/internal/handlers"
)

func New() *http.ServeMux {
	mux := http.NewServeMux()

	mux.HandleFunc("/ping", handlers.Ping)

	mux.HandleFunc("/match/", handlers.FeederByMatch)

	fs := http.FileServer(http.Dir("web"))
	mux.Handle("/", fs)

	return mux
}
