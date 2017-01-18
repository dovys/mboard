package services

import (
	"context"
	"time"

	"github.com/pkg/errors"

	"github.com/dovys/mboard/posts"
	"github.com/dovys/mboard/posts/pb"
	"github.com/satori/go.uuid"
	"google.golang.org/grpc"
)

var (
	ErrUnableToParseUUID    = errors.New("Unable to parse uuid.")
	ErrServiceReturnedError = errors.New("Service returned error.")
	ErrServiceUnavailable   = errors.New("Service unavailable.")
)

type PostsService posts.Service

func NewPostsService() PostsService {
	// lacks timeouts, proper error handling when the service disappears, etc.
	conn, err := grpc.Dial("localhost:8081", grpc.WithInsecure())

	if err != nil {
		panic(errors.Wrap(ErrServiceUnavailable, err.Error()))
	}

	return &postsService{
		c: pb.NewPostsClient(conn),
	}
}

type postsService struct {
	c pb.PostsClient
}

// messy
func (s *postsService) GetLatestPosts(offset int64, limit int64) ([]*posts.Post, error) {
	r, err := s.c.GetLatestPosts(
		context.Background(),
		&pb.GetLatestPostsRequest{Offset: offset, Limit: limit},
	)

	if err != nil {
		return nil, errors.Wrap(ErrServiceUnavailable, err.Error())
	}

	if r.Err != "" {
		return nil, errors.Wrap(ErrServiceReturnedError, r.Err)
	}

	res := make([]*posts.Post, len(r.Posts))

	for i, p := range r.Posts {
		id, err := uuid.FromBytes(p.Id)

		if err != nil {
			return nil, errors.Wrap(ErrUnableToParseUUID, err.Error())
		}

		res[i] = &posts.Post{
			Id:     id,
			Author: p.Author,
			Text:   p.Text,
			Date:   time.Unix(p.Date, 0),
		}
	}

	return res, nil
}

func (s *postsService) GetPost(id uuid.UUID) (*posts.Post, error) {
	return nil, nil
}

func (s *postsService) AddPost(author string, text string) (uuid.UUID, error) {
	return uuid.Nil, nil
}
