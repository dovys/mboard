package main

import (
	"context"
	"flag"
	"net"
	"os"

	"google.golang.org/grpc"

	"github.com/dovys/mboard/posts"
	"github.com/dovys/mboard/posts/pb"
	"github.com/go-kit/kit/log"
)

func main() {
	var listenAddr = flag.String("grpc.addr", ":9001", "Address for GRPC to listen on.")
	var logger log.Logger
	{
		logger = log.NewLogfmtLogger(os.Stdout)
		logger = log.NewContext(logger).With("ts", log.DefaultTimestampUTC)
		logger = log.NewContext(logger).With("caller", log.DefaultCaller)
	}

	flag.Parse()
	logger.Log("status", "started")

	var service posts.Service
	{
		logger := log.NewContext(logger).WithPrefix("domain", "service").With("service", "Posts")
		service = posts.NewPostsService(logger)
	}

	endpoints := posts.MakeServerEndpoints(service)
	ctx := context.Background()

	// Transport
	{
		tcp, err := net.Listen("tcp", *listenAddr)

		if err != nil {
			logger.Log("err", err)
			return
		}

		server := posts.MakePostsServer(ctx, endpoints, logger)
		s := grpc.NewServer()
		pb.RegisterPostsServer(s, server)

		logger.Log("status", "listening", "addr", *listenAddr)
		s.Serve(tcp)
	}
}
