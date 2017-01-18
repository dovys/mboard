package posts

import (
	"golang.org/x/net/context"

	"github.com/go-kit/kit/endpoint"
	uuid "github.com/satori/go.uuid"
)

type Endpoints struct {
	GetLatestPostsEndpoint endpoint.Endpoint
	GetPostEndpoint        endpoint.Endpoint
	AddPostEndpoint        endpoint.Endpoint
}

func MakeServerEndpoints(s Service) Endpoints {
	return Endpoints{
		GetLatestPostsEndpoint: MakeGetLatestPostsEndpoint(s),
		GetPostEndpoint:        MakeGetPostEndpoint(s),
		AddPostEndpoint:        MakeAddPostEndpoint(s),
	}
}

func MakeGetLatestPostsEndpoint(s Service) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (response interface{}, err error) {
		req := request.(*getLatestPostsRequest)
		posts, err := s.GetLatestPosts(req.Offset, req.Limit)

		return &getLatestPostsResponse{Posts: posts, Err: err}, nil
	}
}

func MakeGetPostEndpoint(s Service) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (response interface{}, err error) {
		req := request.(*getPostRequest)
		post, err := s.GetPost(req.Id)

		return &getPostResponse{
			Post: post,
			Err:  err,
		}, nil
	}
}

func MakeAddPostEndpoint(s Service) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (response interface{}, err error) {
		req := request.(*addPostRequest)
		uuid, err := s.AddPost(req.Author, req.Text)

		return &addPostResponse{Id: uuid, Err: err}, nil
	}
}

// DTO's to use between the transport and the endpoints.
// Makes the endpoints not care about the particular transport impl
type getPostRequest struct {
	Id uuid.UUID
}

type getPostResponse struct {
	Post *Post
	Err  error
}

type addPostRequest struct {
	Author string
	Text   string
}

type addPostResponse struct {
	Id  uuid.UUID
	Err error
}

type getLatestPostsRequest struct {
	Offset int64
	Limit  int64
}

type getLatestPostsResponse struct {
	Posts []*Post
	Err   error
}
