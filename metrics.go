package main

import (
	"fmt"
	"log"
	"net/http"
)

func (cfg *apiConfig) middlewareMetricsInc(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		cfg.fileserverHits.Add(1)
		log.Printf("%s %s", r.Method, r.URL.Path)
		next.ServeHTTP(w, r)
	})
}

func (cfg *apiConfig) handlerMetrics(wr http.ResponseWriter, req *http.Request) {
	wr.Header().Add("Content-Type", "text/html; charset=utf-8")
	wr.WriteHeader(http.StatusOK)
	wr.Write([]byte(fmt.Sprintf(`
	<html>
	<body>
		<h1>Welcome, Chirpy Admin</h1>
		<p>Chirpy has been visited %d times!</p>
	</body>
	</html>
	`, cfg.fileserverHits.Load())))
}
