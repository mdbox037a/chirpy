package main

import (
	_ "github.com/lib/pq"

	"fmt"
	"log"
	"net/http"
	"sync/atomic"
)

const filepathRoot string = "."
const port string = "8080"

type apiConfig struct {
	fileserverHits atomic.Int32
}

func (cfg *apiConfig) middlewareMetricsInc(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		cfg.fileserverHits.Add(1)
		log.Printf("%s %s", r.Method, r.URL.Path)
		next.ServeHTTP(w, r)
	})
}

func main() {
	apiCfg := apiConfig{}

	mux := http.NewServeMux()
	mux.Handle("/app/", apiCfg.middlewareMetricsInc(http.StripPrefix("/app", http.FileServer(http.Dir(filepathRoot)))))
	mux.HandleFunc("GET /api/healthz", handlerReadiness)
	mux.HandleFunc("GET /admin/metrics", apiCfg.handlerMetrics)
	mux.HandleFunc("POST /admin/reset", apiCfg.handlerReset)
	mux.HandleFunc("POST /api/validate_chirp", handlerValidateChirp)

	srv := &http.Server{
		Addr:    ":" + port,
		Handler: mux,
	}

	log.Printf("Serving files from %s on port %s\n", filepathRoot, port)
	err := srv.ListenAndServe()
	if err != nil {
		log.Fatal(err)
	}
}

func handlerReadiness(wr http.ResponseWriter, req *http.Request) {
	wr.Header().Add("Content-Type", "text/plain; charset=utf-8")
	wr.WriteHeader(200)
	wr.Write([]byte("OK\n"))
}

func (cfg *apiConfig) handlerMetrics(wr http.ResponseWriter, req *http.Request) {
	wr.Header().Add("Content-Type", "text/html; charset=utf-8")
	wr.WriteHeader(200)
	wr.Write([]byte(fmt.Sprintf(`
	<html>
	<body>
		<h1>Welcome, Chirpy Admin</h1>
		<p>Chirpy has been visited %d times!</p>
	</body>
	</html>
	`, cfg.fileserverHits.Load())))
}

func (cfg *apiConfig) handlerReset(wr http.ResponseWriter, req *http.Request) {
	previous := cfg.fileserverHits.Load()
	cfg.fileserverHits.Store(0)
	wr.Header().Add("Content-Type", "text/plain; charset=utf-8")
	wr.WriteHeader(200)
	wr.Write([]byte(fmt.Sprintf("Info: reset site hits counter to 0 \n[previous value: %d]\n[current value: %d]\n", previous, cfg.fileserverHits.Load())))
}
