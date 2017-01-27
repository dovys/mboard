package handlers

import (
	"testing"

	"net/http"
	"net/http/httptest"

	"strings"

	"github.com/dovys/mboard/api/handlers/mock"
	"github.com/dovys/mboard/posts"
	"github.com/gorilla/mux"
	uuid "github.com/satori/go.uuid"
	"github.com/stretchr/testify/assert"
)

type PostsHandlerTestSuite struct {
	postsService *mock.MockPostsService
	mux          *mux.Router
	w            *httptest.ResponseRecorder
}

func setupSuite() *PostsHandlerTestSuite {
	s := &PostsHandlerTestSuite{
		postsService: new(mock.MockPostsService),
		mux:          mux.NewRouter(),
		w:            httptest.NewRecorder(),
	}

	NewPostsHandler(s.postsService).Register(s.mux)

	return s
}

func TestGetLatestPosts(t *testing.T) {
	// query -> expected offset & limit in the service call
	cases := map[string][2]int64{
		"":                 [2]int64{0, 10},
		"limit=2&offset=3": [2]int64{3, 2},
		"limit=0":          [2]int64{0, 10},
		"limit=-1":         [2]int64{0, 10},
		"limit=eleven":     [2]int64{0, 10},
		"offset=eleven":    [2]int64{0, 10},
		"offset=20":        [2]int64{20, 10},
		"offset=0":         [2]int64{0, 10},
		"offset=-1":        [2]int64{0, 10},
	}

	for query, with := range cases {
		t.Run(query, func(t *testing.T) {
			s := setupSuite()
			s.postsService.
				On("GetLatestPosts", with[0], with[1]).
				Return([]*posts.Post{&posts.Post{Author: "Suite", Text: "Testing"}}).
				Once()

			// Execute a request, see that the services.PostsService.GetLatestPosts gets called
			// with the expected values and serializes the result of GetLatestPosts as a response.
			s.mux.ServeHTTP(s.w, httptest.NewRequest("GET", "/?"+query, nil))

			assert.Equal(t, http.StatusOK, s.w.Code)
			assert.Equal(t,
				`[{"Id":"00000000-0000-0000-0000-000000000000","Author":"Suite","Text":"Testing","Date":"0001-01-01T00:00:00Z"}]`+"\n",
				s.w.Body.String(),
			)

			s.postsService.AssertExpectations(t)
		})
	}
}

func TestGetPost(t *testing.T) {
	s := setupSuite()

	id, _ := uuid.FromString("5608f3e1-bfc7-495f-9e09-709fefb28dc9")
	s.postsService.
		On("GetPost", id).
		Return(&posts.Post{Id: id}).
		Once()

	s.mux.ServeHTTP(s.w, httptest.NewRequest("GET", "/5608f3e1-bfc7-495f-9e09-709fefb28dc9", nil))

	assert.Equal(t, http.StatusOK, s.w.Code)
	s.postsService.AssertExpectations(t)
}

func TestGetPostInvalidUuid(t *testing.T) {
	s := setupSuite()

	s.mux.ServeHTTP(s.w, httptest.NewRequest("GET", "/invalid", nil))

	assert.Equal(t, http.StatusBadRequest, s.w.Code)
	assert.Equal(t, `{"error":"Invalid UUID."}`+"\n", s.w.Body.String())
	s.postsService.AssertNotCalled(t, "GetPost")
}

func TestGetPostNotFound(t *testing.T) {
	s := setupSuite()

	var nilPost *posts.Post
	id, _ := uuid.FromString("5608f3e1-bfc7-495f-9e09-709fefb28dc9")
	s.postsService.
		On("GetPost", id).
		Return(nilPost).
		Once()

	s.mux.ServeHTTP(s.w, httptest.NewRequest("GET", "/5608f3e1-bfc7-495f-9e09-709fefb28dc9", nil))

	assert.Equal(t, http.StatusNotFound, s.w.Code)
	assert.Equal(t, `{"error":"Post not found."}`+"\n", s.w.Body.String())
	s.postsService.AssertExpectations(t)
}

func TestAddPost(t *testing.T) {
	s := setupSuite()
	id := uuid.NewV4()

	s.postsService.
		On("AddPost", "Suite", "Lowercase").
		Return(id).
		Once()

	body := strings.NewReader(`{"Author":"Suite","text":"Lowercase"}`)
	s.mux.ServeHTTP(s.w, httptest.NewRequest("POST", "/", body))

	assert.Equal(t, http.StatusOK, s.w.Code)
	assert.Equal(t, `"`+id.String()+`"`+"\n", s.w.Body.String())
	s.postsService.AssertExpectations(t)
}

func TestAddPostWithInvalidData(t *testing.T) {
	cases := map[string]string{
		"InvalidJson":      "invalid",
		"Empty":            "",
		"Collection":       "[]",
		"InvalidStructure": `{"Text":"LaLaLaLand","Author":0}`,
	}

	for title, body := range cases {
		t.Run(title, func(t *testing.T) {
			s := setupSuite()
			s.mux.ServeHTTP(s.w, httptest.NewRequest("POST", "/", strings.NewReader(body)))
			assert.Equal(t, http.StatusBadRequest, s.w.Code)
			assert.Equal(t,
				`{"error":"Invalid post data. Expected body: {\"Author\":\"John Smith\",\"Text\":\"Post\"}."}`+"\n",
				s.w.Body.String(),
			)
			s.postsService.AssertNotCalled(t, "AddPost")
		})
	}
}
