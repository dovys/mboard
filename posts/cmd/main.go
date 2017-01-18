package main

import (
	"context"
	"net"
	"os"

	"google.golang.org/grpc"

	"github.com/dovys/mboard/posts"
	"github.com/dovys/mboard/posts/pb"
	"github.com/go-kit/kit/log"
)

func main() {
	var logger log.Logger
	{
		logger = log.NewLogfmtLogger(os.Stdout)
		logger = log.NewContext(logger).With("ts", log.DefaultTimestampUTC)
		logger = log.NewContext(logger).With("caller", log.DefaultCaller)
	}

	logger.Log("msg", "hello", "key", "value")

	var service posts.Service
	{
		logger := log.NewContext(logger).WithPrefix("domain", "service").With("service", "Posts")
		service = posts.NewPostsService(logger)
	}

	endpoints := posts.MakeServerEndpoints(service)
	ctx := context.Background()

	// Transport
	{
		tcp, err := net.Listen("tcp", ":8081")
		if err != nil {
			panic(err)
		}

		server := posts.MakePostsServer(ctx, endpoints, logger)
		s := grpc.NewServer()
		pb.RegisterPostsServer(s, server)

		logger.Log("addr", ":8081", "transport", "grpc")
		s.Serve(tcp)
	}
}
