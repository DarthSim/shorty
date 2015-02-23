package main

import (
	"net"
	"net/http"

	"github.com/gorilla/mux"
	"github.com/rcrowley/goagain"
)

// Server tools ====================================================================================

func initServer() {
	listener, err := goagain.Listener()

	if err != nil {
		bindAddress := config.Server.Address + ":" + config.Server.Port

		logger.Printf("Starting server on %s\n", bindAddress)

		listener, err = net.Listen("tcp", bindAddress)
		if err != nil {
			logger.Fatalf("Can't start server: %v", err)
		}

		go startServer(listener)
	} else {
		logger.Printf("Resume server on %s\n", listener.Addr())

		go startServer(listener)

		if err := goagain.Kill(); nil != err {
			logger.Fatalf("Can't resume server: %v", err)
		}
	}

	if _, err := goagain.Wait(listener); nil != err {
		logger.Fatalln(err)
	}
}

func startServer(l net.Listener) {
	if err := http.Serve(l, setupRouter()); err != nil {
		logger.Fatalf("Can't start server: %v", err)
	}
}

func setupRouter() (router *mux.Router) {
	router = mux.NewRouter()

	router.HandleFunc("/shorten", createUrlHandler).Methods("POST")
	router.HandleFunc("/{code}", redirectHandler).Methods("GET")
	router.HandleFunc("/expand/{code}", expandHandler).Methods("GET")
	router.HandleFunc("/statistics/{code}", statisticsHandler).Methods("GET")

	return
}

func requestVars(req *http.Request) map[string]string {
	return mux.Vars(req)
}

func serverError(rw http.ResponseWriter, err error, status int) {
	logger.Printf("Server error: %v", err)

	rw.Header().Set("Content-Type", "text/plain; charset=utf-8")
	rw.WriteHeader(status)
	rw.Write([]byte("Internal server error"))
}

func serverResponse(rw http.ResponseWriter, response string, status int) {
	rw.Header().Set("Content-Type", "text/plain; charset=utf-8")
	rw.WriteHeader(status)
	rw.Write([]byte(response))
}

// end of Server tools
