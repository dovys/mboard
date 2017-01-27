package mock

import (
	"github.com/dovys/mboard/posts"
	uuid "github.com/satori/go.uuid"
	"github.com/stretchr/testify/mock"
)

type MockPostsService struct {
	mock.Mock
}

func (m *MockPostsService) GetLatestPosts(offset int64, limit int64) ([]*posts.Post, error) {
	args := m.Called(offset, limit)

	return args.Get(0).([]*posts.Post), nil
}

func (m *MockPostsService) GetPost(id uuid.UUID) (*posts.Post, error) {
	args := m.Called(id)

	return args.Get(0).(*posts.Post), nil
}

func (m *MockPostsService) AddPost(author string, text string) (uuid.UUID, error) {
	args := m.Called(author, text)

	return args.Get(0).(uuid.UUID), nil
}
