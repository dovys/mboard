package main

import (
	"flag"
	"log"
	"net/http"
	"os"

	"github.com/dovys/mboard/handlers"
	"github.com/dovys/mboard/services"
	"github.com/gorilla/mux"
)

func main() {
	address := *flag.String("host", "localhost:8080", "Host.")

	hostname, _ := os.Hostname()
	mux := mux.NewRouter()

	mux.HandleFunc("/status", func(rw http.ResponseWriter, r *http.Request) {
		rw.WriteHeader(http.StatusOK)
		rw.Write([]byte("Running on "))
		rw.Write([]byte(hostname))
	})

	handlers.NewPostsHandler(services.NewPostsService()).Register(mux.PathPrefix("/posts").Subrouter())

	mux.Handle("/", http.FileServer(http.Dir("./static/")))

	// _ "net/http/pprof"
	// go func() {
	// 	log.Println(http.ListenAndServe(":6060", nil))
	// }()

	log.Printf("Running %s on %s\n", address, hostname)
	log.Fatal(http.ListenAndServe(address, mux))
}
