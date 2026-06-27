package main

import (
	"fmt"
	"net/http"
)

func (cfg *apiConfig) handlerReset(wr http.ResponseWriter, req *http.Request) {
	previous := cfg.fileserverHits.Load()
	cfg.fileserverHits.Store(0)
	wr.Header().Add("Content-Type", "text/plain; charset=utf-8")
	wr.WriteHeader(200)
	wr.Write([]byte(fmt.Sprintf("Info: reset site hits counter to 0 \n[previous value: %d]\n[current value: %d]\n", previous, cfg.fileserverHits.Load())))
}
