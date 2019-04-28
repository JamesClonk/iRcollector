package main

import (
	"net/http"

	"github.com/JamesClonk/iRcollector/env"
	"github.com/JamesClonk/iRcollector/log"
	"github.com/gorilla/mux"
)

func main() {
	port := env.Get("PORT", "8080")
	level := env.Get("LOG_LEVEL", "info")

	log.Infoln("port:", port)
	log.Infoln("log level:", level)

	// start listener
	log.Fatalln(http.ListenAndServe(":"+port, router()))
}

func router() *mux.Router {
	r := mux.NewRouter()
	r.PathPrefix("/health").HandlerFunc(health)
	return r
}

func health(rw http.ResponseWriter, req *http.Request) {
	rw.WriteHeader(200)
	rw.Header().Set("X-IRcollector", "IRcollector")
	rw.Header().Set("Content-Type", "application/json")
	rw.Write([]byte(`{ "status": "ok" }`))
}
