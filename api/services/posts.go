package services

import (
	"context"
	"time"

	"github.com/go-kit/kit/log"

	"github.com/pkg/errors"

	"github.com/dovys/mboard/posts"
	"github.com/dovys/mboard/posts/pb"
	"github.com/satori/go.uuid"
)

var (
	ErrServiceUnavailable = errors.New("Service unavailable.")
)

type PostsService posts.Service

func NewPostsService(client pb.PostsClient, logger log.Logger) PostsService {
	return &postsService{
		client: client,
		logger: logger,
	}
}

type postsService struct {
	client pb.PostsClient
	logger log.Logger
}

func (s *postsService) GetLatestPosts(offset int64, limit int64) ([]*posts.Post, error) {
	s.logger.Log("limit", limit, "offset", offset)
	res, err := s.client.GetLatestPosts(
		context.Background(),
		&pb.GetLatestPostsRequest{Offset: offset, Limit: limit},
	)

	if err != nil {
		s.logger.Log("err", err)
		return nil, ErrServiceUnavailable
	}

	if res.Err != "" {
		return nil, errors.New(res.Err)
	}

	collection := make([]*posts.Post, len(res.Posts))

	for i, p := range res.Posts {
		collection[i] = s.decodePost(p)
	}

	return collection, nil
}

func (s *postsService) GetPost(id uuid.UUID) (*posts.Post, error) {
	res, err := s.client.GetPost(
		context.Background(),
		&pb.GetPostRequest{Id: id.Bytes()},
	)

	if err != nil {
		s.logger.Log("err", err)
		return nil, ErrServiceUnavailable
	}

	if res.Err != "" {
		return nil, errors.New(res.Err)
	}

	// 404
	if len(res.Post.Id) == 0 {
		return nil, nil
	}

	return s.decodePost(res.Post), nil
}

func (s *postsService) AddPost(author string, text string) (uuid.UUID, error) {
	res, err := s.client.AddPost(
		context.Background(),
		&pb.AddPostRequest{Author: author, Text: text},
	)

	if err != nil {
		s.logger.Log("err", err)
		return uuid.Nil, ErrServiceUnavailable
	}

	if res.Err != "" {
		return uuid.Nil, errors.New(res.Err)
	}

	id, err := uuid.FromBytes(res.Id)

	if err != nil {
		s.logger.Log("err", err)
		return uuid.Nil, ErrServiceUnavailable
	}

	return id, nil
}

func (s *postsService) decodePost(p *pb.Post) *posts.Post {
	id, err := uuid.FromBytes(p.Id)

	if err != nil {
		s.logger.Log("err", err)
		return nil
	}

	return &posts.Post{
		Id:     id,
		Author: p.Author,
		Text:   p.Text,
		Date:   time.Unix(p.Date, 0),
	}
}
