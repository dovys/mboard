package posts

import (
	"github.com/go-kit/kit/log"

	"github.com/dovys/mboard/posts/pb"

	grpctransport "github.com/go-kit/kit/transport/grpc"
	uuid "github.com/satori/go.uuid"
	"golang.org/x/net/context"
)

func MakePostsServer(ctx context.Context, endpoints Endpoints, logger log.Logger) pb.PostsServer {
	options := []grpctransport.ServerOption{
		grpctransport.ServerErrorLogger(logger),
	}

	return &postsServer{
		getLatestPostsHandler: grpctransport.NewServer(
			ctx,
			endpoints.GetLatestPostsEndpoint,
			DecodeGetLatestPostsRequest,
			EncodeGetLatestPostsResponse,
			options...,
		),
		getPostHandler: grpctransport.NewServer(
			ctx,
			endpoints.GetPostEndpoint,
			DecodeGetPostRequest,
			EncodeGetPostResponse,
			options...,
		),
		addPostHandler: grpctransport.NewServer(
			ctx,
			endpoints.AddPostEndpoint,
			DecodeAddPostRequest,
			EncodeAddPostResponse,
			options...,
		),
	}
}

// Implements pb.PostsServer
type postsServer struct {
	getLatestPostsHandler grpctransport.Handler
	getPostHandler        grpctransport.Handler
	addPostHandler        grpctransport.Handler
}

func (s *postsServer) GetPost(ctx context.Context, req *pb.GetPostRequest) (*pb.GetPostResponse, error) {
	_, response, err := s.getPostHandler.ServeGRPC(ctx, req)

	if err != nil {
		return nil, err
	}
	return response.(*pb.GetPostResponse), nil
}

func (s *postsServer) GetLatestPosts(ctx context.Context, req *pb.GetLatestPostsRequest) (*pb.GetLatestPostsResponse, error) {
	_, response, err := s.getLatestPostsHandler.ServeGRPC(ctx, req)

	if err != nil {
		return nil, err
	}

	return response.(*pb.GetLatestPostsResponse), nil
}

func (s *postsServer) AddPost(ctx context.Context, req *pb.AddPostRequest) (*pb.AddPostResponse, error) {
	_, response, err := s.addPostHandler.ServeGRPC(ctx, req)

	if err != nil {
		return nil, err
	}

	return response.(*pb.AddPostResponse), nil
}

// todo: Don't export
func DecodeGetLatestPostsRequest(_ context.Context, grpcRequest interface{}) (interface{}, error) {
	req := grpcRequest.(*pb.GetLatestPostsRequest)

	return &getLatestPostsRequest{Offset: req.Offset, Limit: req.Limit}, nil
}

func EncodeGetLatestPostsResponse(_ context.Context, endpointResponse interface{}) (interface{}, error) {
	response := endpointResponse.(*getLatestPostsResponse)
	posts := make([]*pb.Post, len(response.Posts))

	if response.Err != nil {
		return &pb.GetLatestPostsResponse{Posts: posts, Err: response.Err.Error()}, nil
	}

	for i, post := range response.Posts {
		posts[i] = &pb.Post{
			Id:     post.Id.Bytes(),
			Author: post.Author,
			Text:   post.Text,
			Date:   post.Date.Unix(),
		}
	}

	return &pb.GetLatestPostsResponse{Posts: posts}, nil
}

func DecodeGetPostRequest(_ context.Context, grpcRequest interface{}) (interface{}, error) {
	req := grpcRequest.(*pb.GetPostRequest)

	id, err := uuid.FromBytes(req.Id)

	if err != nil || id == uuid.Nil {
		return nil, err
	}

	return &getPostRequest{Id: id}, nil
}

func EncodeGetPostResponse(_ context.Context, endpointResponse interface{}) (interface{}, error) {
	res := endpointResponse.(*getPostResponse)

	if res.Post == nil || res.Err != nil {
		var err string
		if res.Err != nil {
			err = res.Err.Error()
		}

		return &pb.GetPostResponse{Post: &pb.Post{}, Err: err}, nil
	}

	return &pb.GetPostResponse{
		Post: &pb.Post{
			Id:     res.Post.Id.Bytes(),
			Author: res.Post.Author,
			Text:   res.Post.Text,
			Date:   res.Post.Date.Unix(),
		},
	}, nil
}

func DecodeAddPostRequest(_ context.Context, grpcRequest interface{}) (interface{}, error) {
	req := grpcRequest.(*pb.AddPostRequest)

	return &addPostRequest{Author: req.Author, Text: req.Text}, nil
}

func EncodeAddPostResponse(_ context.Context, endpointResponse interface{}) (interface{}, error) {
	res := endpointResponse.(*addPostResponse)

	if res.Err != nil {
		return &pb.AddPostResponse{Err: res.Err.Error()}, nil
	}

	return &pb.AddPostResponse{Id: res.Id.Bytes()}, nil
}
