package services

import (
	"time"

	"sync"

	"github.com/satori/go.uuid"
)

type Post struct {
	Id     uuid.UUID
	Author string
	Text   string
	Date   time.Time
}

type PostsService interface {
	GetPosts(limit uint64, offset uint64) []*Post
	GetPost(uuid.UUID) *Post
	SubmitPost(author string, text string) *uuid.UUID
}

func NewPostsService() PostsService {
	return &postsService{
		posts:       make([]*Post, 0),
		postMap:     make(map[uuid.UUID]*Post, 0),
		submitMutex: &sync.Mutex{},
	}
}

type postsService struct {
	posts       []*Post
	postMap     map[uuid.UUID]*Post
	submitMutex *sync.Mutex
}

func (s *postsService) GetPosts(offset uint64, limit uint64) []*Post {
	if int(offset) >= len(s.posts) {
		return s.posts[0:0]
	}

	if int(offset+limit) > len(s.posts) {
		limit = uint64(len(s.posts)) - offset
	}

	return s.posts[offset : limit+offset]
}

func (s *postsService) GetPost(id uuid.UUID) *Post {
	return s.postMap[id]
}

func (s *postsService) SubmitPost(author string, text string) *uuid.UUID {
	post := &Post{
		Id:     uuid.NewV4(),
		Author: author,
		Text:   text,
		Date:   time.Now(),
	}

	s.submitMutex.Lock()
	defer s.submitMutex.Unlock()

	s.posts = append([]*Post{post}, s.posts...)
	s.postMap[post.Id] = post

	return &post.Id
}
