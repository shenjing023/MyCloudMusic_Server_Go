package main

import (
	"net/http"
	"time"
)

func main() {

	router := NewRouter(allRoutes())

	server := &http.Server{
		Addr:         ":8000",
		ReadTimeout:  60 * time.Second,
		WriteTimeout: 60 * time.Second,
		Handler:      router,
	}
	server.ListenAndServe()
}
