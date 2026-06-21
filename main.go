package main

import (
	"fmt"
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
		fmt.Printf("%v\n", err)
	}
}
