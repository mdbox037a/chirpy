package main

import (
	"log"
	"net/http"
)

func main() {
	srvMux := http.NewServeMux()
	svr := http.Server{
		Addr:    ":8080",
		Handler: srvMux,
	}
	err := svr.ListenAndServe()
	if err != nil {
		log.Fatal(err)
	}
}
