package main

import (
	"context"
	"dotaWorstPlayerChacker/internal/handlers"
	"dotaWorstPlayerChacker/internal/openDota"
	redisstore "dotaWorstPlayerChacker/internal/redisStore"
	"log"
	"net/http"
	"time"

	"dotaWorstPlayerChacker/internal/router"
)

func main() {
	mux := router.New()

	od := openDota.NewClient()
	store := redisstore.NewRedis()

	handlers.SetDeps(handlers.Deps{
		OD:    od,
		Store: store,
	})

	heroes, err := od.GetHeroes()
	if err != nil {
		log.Printf("не вдалося підтягнути героїв на старті: %v", err)
	} else if err := store.RefreshHeroes(context.Background(), heroes); err != nil {
		log.Printf("не вдалося оновити героїв у Redis: %v", err)
	}

	srv := &http.Server{
		Addr:              ":8080",
		Handler:           mux,
		ReadTimeout:       5 * time.Second,
		ReadHeaderTimeout: 10 * time.Second,
		WriteTimeout:      15 * time.Second,
		IdleTimeout:       60 * time.Second,
	}

	log.Println("Starting server on port 8080")
	log.Fatal(srv.ListenAndServe())
}
