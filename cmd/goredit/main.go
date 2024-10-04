package main

import (
	"log"
	"net/http"

	"github.com/alaalser/goreddit/postgres"
	"github.com/alaalser/goreddit/web"
)

func main() {
	dsn := "postgres://postgres:secret@localhost/postgres?sslmode=disable"

	store, err := postgres.NewStore(dsn)
	if err != nil {
		log.Fatal(err)
	}

	h := web.NewHandler(store)
	http.ListenAndServe(":3000", h)
}
