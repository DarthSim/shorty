package main

import (
	"log"
	"net/http"

	"github.com/gorilla/mux"
	"github.com/peterhellberg/env"
)

// Server tools ====================================================================================

func startServer() {
	bindAddress := env.String("ADDRESS", ":8080")

	log.Printf("Starting server on %s\n", bindAddress)

	if err := http.ListenAndServe(bindAddress, setupRouter()); err != nil {
		log.Fatalf("Can't start server: %v", err)
	}
}

func setupRouter() (router *mux.Router) {
	router = mux.NewRouter()

	router.HandleFunc("/shorten", createUrlHandler).Methods("POST")
	router.HandleFunc("/{code}", redirectHandler).Methods("GET", "HEAD")
	router.HandleFunc("/expand/{code}", expandHandler).Methods("GET", "HEAD")
	router.HandleFunc("/statistics/{code}", statisticsHandler).Methods("GET", "HEAD")

	return
}

func serverError(rw http.ResponseWriter, req *http.Request, err interface{}, status int) {
	log.Printf("Server error: %v", err)
	serverResponse(rw, req, "Internal server error", status)
}

func serverResponse(rw http.ResponseWriter, req *http.Request, response string, status int) {
	rw.Header().Set("Content-Type", "text/plain; charset=utf-8")
	rw.WriteHeader(status)

	if req.Method != "HEAD" {
		rw.Write([]byte(response))
	}
}

// end of Server tools
