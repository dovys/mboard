package main

import (
	"flag"
	"net/http"
	"os"

	"google.golang.org/grpc"

	"github.com/go-kit/kit/log"

	"github.com/dovys/mboard/api/handlers"
	"github.com/dovys/mboard/api/services"
	"github.com/dovys/mboard/posts/pb"
	"github.com/gorilla/mux"
)

func main() {
	var (
		address   = flag.String("api.addr", ":8080", "Address to expose API under.")
		apiPrefix = flag.String("api.prefix", "/api", "API Prefix")
	)

	var logger log.Logger
	{
		logger = log.NewLogfmtLogger(os.Stdout)
		logger = log.NewContext(logger).With("ts", log.DefaultTimestampUTC)
		logger = log.NewContext(logger).With("caller", log.DefaultCaller)
	}

	flag.Parse()
	m := mux.NewRouter()
	api := m.PathPrefix(*apiPrefix).Subrouter()

	// /posts Handler
	{
		// lacks timeouts, proper error handling when the service disappears, etc.
		conn, err := grpc.Dial("posts:9001", grpc.WithInsecure())

		if err != nil {
			logger.Log("err", err)
			return
		}

		s := services.NewPostsService(
			pb.NewPostsClient(conn),
			log.NewContext(logger).With("service", "posts"),
		)

		h := handlers.NewPostsHandler(s)
		h.Register(api.PathPrefix("/posts").Subrouter())
	}

	hostname, _ := os.Hostname()
	// Debug
	{
		m.HandleFunc("/status", func(rw http.ResponseWriter, r *http.Request) {
			rw.WriteHeader(http.StatusOK)
			rw.Write([]byte("Running on "))
			rw.Write([]byte(hostname))
		})
	}

	logger.Log("listening", *address, "api", *apiPrefix, "machine", hostname)
	logger.Log("err", http.ListenAndServe(*address, m))
}
