package main

import (
	"net/http"
)

func handlerReadiness(wr http.ResponseWriter, req *http.Request) {
	wr.Header().Add("Content-Type", "text/plain; charset=utf-8")
	wr.WriteHeader(http.StatusOK)
	wr.Write([]byte("OK\n"))
}
