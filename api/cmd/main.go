package main

import (
	"flag"
	"net/http"
	"os"

	"github.com/go-kit/kit/log"

	"net/http/pprof"

	"github.com/dovys/mboard/api/handlers"
	"github.com/dovys/mboard/api/services"
	"github.com/gorilla/mux"
)

func main() {
	var logger log.Logger
	{
		logger = log.NewLogfmtLogger(os.Stdout)
		logger = log.NewContext(logger).With("ts", log.DefaultTimestampUTC)
		logger = log.NewContext(logger).With("caller", log.DefaultCaller)
	}

	address := flag.String("host", ":8080", "Host.")
	hostname, _ := os.Hostname()
	mux := mux.NewRouter()

	mux.HandleFunc("/status", func(rw http.ResponseWriter, r *http.Request) {
		rw.WriteHeader(http.StatusOK)
		rw.Write([]byte("Running on "))
		rw.Write([]byte(hostname))
	})

	// PostsService
	{
		logger := log.NewContext(logger).With("handler", "posts")
		h := handlers.NewPostsHandler(services.NewPostsService(), logger)
		h.Register(mux.PathPrefix("/posts").Subrouter())
	}

	mux.Handle("/", http.FileServer(http.Dir("./static/")))

	logger.Log("listening", *address, "machine", hostname)

	go logger.Log("err", http.ListenAndServe(*address, mux))

	// pprof
	{
		m := http.NewServeMux()
		m.Handle("/debug/pprof/", http.HandlerFunc(pprof.Index))
		m.Handle("/debug/pprof/cmdline", http.HandlerFunc(pprof.Cmdline))
		m.Handle("/debug/pprof/profile", http.HandlerFunc(pprof.Profile))
		m.Handle("/debug/pprof/symbol", http.HandlerFunc(pprof.Symbol))
		m.Handle("/debug/pprof/trace", http.HandlerFunc(pprof.Trace))
		http.ListenAndServe(":33377", nil)
	}
}
