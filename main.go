package main

import (
	"log"
	"net/http"
)

func main() {
	mux := http.NewServeMux()
	mux.Handle("/", http.FileServer(http.Dir(".")))

	svr := http.Server{
		Addr:    ":8080",
		Handler: mux,
	}
	err := svr.ListenAndServe()
	if err != nil {
		log.Fatal(err)
	}
}
