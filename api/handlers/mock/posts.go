package mock

import (
	"github.com/dovys/mboard/api/services"
	uuid "github.com/satori/go.uuid"
	"github.com/stretchr/testify/mock"
)

type MockPostsService struct {
	mock.Mock
}

func (m *MockPostsService) GetLatestPosts(offset uint64, limit uint64) []*services.Post {
	args := m.Called(offset, limit)

	return args.Get(0).([]*services.Post)
}

func (m *MockPostsService) GetPost(id uuid.UUID) *services.Post {
	args := m.Called(id)

	return args.Get(0).(*services.Post)
}

func (m *MockPostsService) AddPost(author string, text string) *uuid.UUID {
	args := m.Called(author, text)

	return args.Get(0).(*uuid.UUID)
}
