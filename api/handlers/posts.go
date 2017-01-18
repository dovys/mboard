package handlers

import (
	"encoding/json"
	"net/http"

	"github.com/go-kit/kit/log"

	"strconv"

	"github.com/dovys/mboard/api/services"
	"github.com/gorilla/mux"
	uuid "github.com/satori/go.uuid"
)

type postsHandler struct {
	postsService services.PostsService
	logger       log.Logger
}

func NewPostsHandler(s services.PostsService, l log.Logger) Handler {
	return &postsHandler{s, l}
}

func (h *postsHandler) handleGetLatestPosts(rw http.ResponseWriter, r *http.Request) {
	var limit, offset int64
	var err error
	e := json.NewEncoder(rw)

	// the max limit will be 2^7-1=127
	if limit, err = strconv.ParseInt(r.URL.Query().Get("limit"), 10, 7); err != nil || limit == 0 {
		limit = 10
	}

	if offset, err = strconv.ParseInt(r.URL.Query().Get("offset"), 10, 64); err != nil {
		offset = 0
	}

	collection, err := h.postsService.GetLatestPosts(offset, limit)

	if err != nil {
		h.logger.Log("method", "GetLatestPosts", "err", err.Error())
		rw.WriteHeader(http.StatusInternalServerError)
		e.Encode(ErrorMessage{"Could not retrieve posts."})
		return
	}

	e.Encode(collection)
}

func (h *postsHandler) handleAddPost(rw http.ResponseWriter, r *http.Request) {
	d := json.NewDecoder(r.Body)
	e := json.NewEncoder(rw)
	defer r.Body.Close()

	data := &struct {
		Author string
		Text   string
	}{}

	if err := d.Decode(data); err != nil {
		h.logger.Log("method", "AddPost", "err", err.Error())
		rw.WriteHeader(http.StatusBadRequest)
		e.Encode(ErrorMessage{"Invalid post data. Expected body: {\"Author\":\"John Smith\",\"Text\":\"Post\"}."})
		return
	}

	id, err := h.postsService.AddPost(data.Author, data.Text)

	if err != nil {
		h.logger.Log("method", "AddPost", "err", err.Error())
		// todo different headers for different errors?
		rw.WriteHeader(http.StatusServiceUnavailable)
		e.Encode(ErrorMessage{err.Error()})
		return
	}

	e.Encode(id)
}

func (h *postsHandler) handleGetPost(rw http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	e := json.NewEncoder(rw)
	id, err := uuid.FromString(vars["post"])

	if err != nil {
		h.logger.Log("method", "GetPost", "err", err.Error())
		rw.WriteHeader(http.StatusBadRequest)
		e.Encode(ErrorMessage{"Invalid UUID."})
		return
	}

	post, err := h.postsService.GetPost(id)

	if err != nil {
		h.logger.Log("method", "GetPost", "err", err.Error())
		// todo different headers for different errors?
		rw.WriteHeader(http.StatusNotFound)
		e.Encode(ErrorMessage{"Post not found."})
		return
	}

	e.Encode(post)
}

func (h *postsHandler) Register(mux *mux.Router) {
	mux.HandleFunc("/", h.handleGetLatestPosts).Methods("GET")
	mux.HandleFunc("/", h.handleAddPost).Methods("POST")
	mux.HandleFunc("/{post}", h.handleGetPost).Methods("GET")
}
