package handler

import (
	"context"
	"encoding/json"
	"log"
	"net/http"
	"strconv"

	"github.com/romakorinenko/hw-test/hw15_go_sql/internal/repository"
)

type IUserHandler interface {
	Handle(w http.ResponseWriter, r *http.Request)
	GetAll(w http.ResponseWriter, r *http.Request)
}

type UserHandler struct {
	userRepository repository.IUserRepository
}

func NewUserHandler(userRepository repository.IUserRepository) IUserHandler {
	return &UserHandler{userRepository: userRepository}
}

func (h *UserHandler) Handle(w http.ResponseWriter, r *http.Request) {
	log.Println("received request: ", r.Method, r.URL.Path)

	var response *Response
	var err error

	switch r.Method {
	case http.MethodPost:
		response, err = h.create(r)
	case http.MethodGet:
		response, err = h.getByID(r)
	case http.MethodPut:
		response, err = h.update(r)
	case http.MethodDelete:
		response, err = h.deleteByID(r)
	}

	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(response.StatusCode)
	if _, err = w.Write(response.Body); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
}

func (h *UserHandler) create(r *http.Request) (*Response, error) {
	var user repository.User
	err := json.NewDecoder(r.Body).Decode(&user)
	if err != nil {
		return nil, err
	}

	createdUser, err := h.userRepository.Create(context.Background(), &user)
	if err != nil {
		return nil, err
	}

	userJson, err := json.Marshal(createdUser)
	if err != nil {
		return nil, err
	}

	return &Response{
		StatusCode: http.StatusCreated,
		Body:       userJson,
	}, nil
}

func (h *UserHandler) update(r *http.Request) (*Response, error) {
	var user repository.User
	err := json.NewDecoder(r.Body).Decode(&user)
	if err != nil {
		return nil, err
	}

	updatedUser, err := h.userRepository.Update(context.Background(), &user)
	if err != nil {
		return nil, err
	}

	userJson, err := json.Marshal(updatedUser)
	if err != nil {
		return nil, err
	}

	return &Response{
		StatusCode: http.StatusOK,
		Body:       userJson,
	}, nil
}

func (h *UserHandler) deleteByID(r *http.Request) (*Response, error) {
	userIdString := r.URL.Query().Get("id")
	userId, err := strconv.Atoi(userIdString)
	if err != nil {
		return nil, err
	}

	err = h.userRepository.DeleteByID(context.Background(), userId)
	if err != nil {
		return nil, err
	}
	return &Response{
		StatusCode: http.StatusOK,
	}, nil
}

func (h *UserHandler) getByID(r *http.Request) (*Response, error) {
	userIdString := r.URL.Query().Get("id")
	userId, err := strconv.Atoi(userIdString)
	if err != nil {
		return nil, err
	}

	user, err := h.userRepository.GetByID(context.Background(), userId)
	userJson, err := json.Marshal(user)
	if err != nil {
		return nil, err
	}

	return &Response{
		StatusCode: http.StatusOK,
		Body:       userJson,
	}, nil
}

func (h *UserHandler) GetAll(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, "", http.StatusMethodNotAllowed)
		return
	}

	users, err := h.userRepository.GetAll(context.Background())
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	userJson, err := json.Marshal(users)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
	_, err = w.Write(userJson)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
}
