package handlers

import (
	"github.com/dovys/mboard/services"
	uuid "github.com/satori/go.uuid"
	"github.com/stretchr/testify/mock"
)

type mockPostsService struct {
	mock.Mock
}

func (m *mockPostsService) GetPosts(limit uint64, offset uint64) []*services.Post {
	args := m.Called(limit, offset)

	return args.Get(0).([]*services.Post)
}

func (m *mockPostsService) GetPost(id uuid.UUID) *services.Post {
	args := m.Called(id)

	return args.Get(0).(*services.Post)
}

func (m *mockPostsService) SubmitPost(author string, text string) *uuid.UUID {
	args := m.Called(author, text)

	return args.Get(0).(*uuid.UUID)
}
