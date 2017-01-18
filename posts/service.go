package posts

import (
	"sync"
	"time"

	"github.com/pkg/errors"

	"github.com/go-kit/kit/log"

	uuid "github.com/satori/go.uuid"
)

type Service interface {
	GetPost(uuid.UUID) (*Post, error)
	GetLatestPosts(offset int64, limit int64) ([]*Post, error)
	AddPost(author string, text string) (uuid.UUID, error)
}

type Post struct {
	Id     uuid.UUID
	Author string
	Text   string
	Date   time.Time
}

var (
	ErrPostNotFound = errors.New("POST_NOT_FOUND")
	ErrPostInvalid  = errors.New("POST_INVALID")
)

func NewPostsService(logger log.Logger) Service {
	return &postsService{
		posts:     make([]*Post, 0),
		postMap:   make(map[uuid.UUID]*Post, 0),
		logger:    logger,
		addLocker: &sync.Mutex{},
	}
}

type postsService struct {
	posts     []*Post
	postMap   map[uuid.UUID]*Post
	logger    log.Logger
	addLocker sync.Locker
}

func (s *postsService) GetLatestPosts(offset int64, limit int64) ([]*Post, error) {
	s.AddPost("user", "post")
	if int(offset) >= len(s.posts) || offset < 0 || limit < 1 {
		return s.posts[0:0], nil
	}

	if int(offset+limit) > len(s.posts) {
		limit = int64(len(s.posts)) - offset
	}

	return s.posts[offset : limit+offset], nil
}

func (s *postsService) GetPost(id uuid.UUID) (*Post, error) {
	if post, ok := s.postMap[id]; ok {
		return post, nil
	}

	log.NewContext(s.logger).With("id", id.String()).Log("err", ErrPostNotFound.Error())

	return nil, ErrPostNotFound
}

func (s *postsService) AddPost(author string, text string) (uuid.UUID, error) {
	if author == "" || text == "" {
		return uuid.Nil, ErrPostInvalid
	}

	post := &Post{
		Id:     uuid.NewV4(),
		Author: author,
		Text:   text,
		Date:   time.Now(),
	}

	s.addLocker.Lock()
	defer s.addLocker.Unlock()

	s.posts = append([]*Post{post}, s.posts...)
	s.postMap[post.Id] = post

	return post.Id, nil
}
