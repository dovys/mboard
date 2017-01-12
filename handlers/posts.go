package handlers

import (
	"encoding/json"
	"net/http"

	"strconv"

	"github.com/dovys/mboard/services"
	"github.com/gorilla/mux"
	uuid "github.com/satori/go.uuid"
)

type postsHandler struct {
	postsService services.PostsService
}

func NewPostsHandler(s services.PostsService) Handler {
	return &postsHandler{s}
}

func (h *postsHandler) handleGetPosts(rw http.ResponseWriter, r *http.Request) {
	var limit, offset uint64
	var err error

	// the max limit will be 2^7-1=127
	if limit, err = strconv.ParseUint(r.URL.Query().Get("limit"), 10, 7); err != nil || limit == 0 {
		limit = 10
	}

	if offset, err = strconv.ParseUint(r.URL.Query().Get("offset"), 10, 64); err != nil {
		offset = 0
	}

	json.NewEncoder(rw).Encode(h.postsService.GetPosts(offset, limit))
}

func (h *postsHandler) handleSubmitPost(rw http.ResponseWriter, r *http.Request) {
	d := json.NewDecoder(r.Body)
	e := json.NewEncoder(rw)
	defer r.Body.Close()

	data := &struct {
		Author string
		Text   string
	}{}

	if err := d.Decode(data); err != nil || data.Author == "" || data.Text == "" {
		rw.WriteHeader(http.StatusBadRequest)
		e.Encode(ErrorMessage{"Invalid post data. Expected body: {\"Author\":\"John Smith\",\"Text\":\"Post\"}."})
		return
	}

	e.Encode(h.postsService.SubmitPost(data.Author, data.Text))
}

func (h *postsHandler) handleGetPost(rw http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	e := json.NewEncoder(rw)
	id, err := uuid.FromString(vars["post"])

	if err != nil {
		rw.WriteHeader(http.StatusBadRequest)
		e.Encode(ErrorMessage{"Invalid UUID."})
		return
	}

	if post := h.postsService.GetPost(id); post != nil {
		e.Encode(post) // @todo: handle the error
	} else {
		rw.WriteHeader(http.StatusNotFound)
		e.Encode(ErrorMessage{"Post not found."})
	}
}

func (h *postsHandler) Register(mux *mux.Router) {
	mux.HandleFunc("/", h.handleGetPosts).Methods("GET")
	mux.HandleFunc("/", h.handleSubmitPost).Methods("POST")
	mux.HandleFunc("/{post}", h.handleGetPost).Methods("GET")
}
