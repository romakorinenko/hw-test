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

func NewUserHandler(userRepository repository.IUserRepository) *UserHandler {
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
	if _, writeBodyErr := w.Write(response.Body); err != nil {
		http.Error(w, writeBodyErr.Error(), http.StatusInternalServerError)
		return
	}
}

func (h *UserHandler) create(r *http.Request) (*Response, error) {
	var user repository.User
	err := json.NewDecoder(r.Body).Decode(&user)
	if err != nil {
		return nil, err
	}

	createdUser, createUserErr := h.userRepository.Create(context.Background(), &user)
	if createUserErr != nil {
		return nil, createUserErr
	}

	userJSON, marshalJSONErr := json.Marshal(createdUser)
	if marshalJSONErr != nil {
		return nil, marshalJSONErr
	}
	return &Response{
		StatusCode: http.StatusCreated,
		Body:       userJSON,
	}, nil
}

func (h *UserHandler) update(r *http.Request) (*Response, error) {
	var user repository.User
	err := json.NewDecoder(r.Body).Decode(&user)
	if err != nil {
		return nil, err
	}

	updatedUser, updateUserErr := h.userRepository.Update(context.Background(), &user)
	if updateUserErr != nil {
		return nil, updateUserErr
	}

	userJSON, marshalJSONErr := json.Marshal(updatedUser)
	if marshalJSONErr != nil {
		return nil, marshalJSONErr
	}
	return &Response{
		StatusCode: http.StatusOK,
		Body:       userJSON,
	}, nil
}

func (h *UserHandler) deleteByID(r *http.Request) (*Response, error) {
	userIDString := r.URL.Query().Get("id")
	userID, err := strconv.Atoi(userIDString)
	if err != nil {
		return nil, err
	}

	deleteUserErr := h.userRepository.DeleteByID(context.Background(), userID)
	if deleteUserErr != nil {
		return nil, deleteUserErr
	}
	return &Response{
		StatusCode: http.StatusOK,
	}, nil
}

func (h *UserHandler) getByID(r *http.Request) (*Response, error) {
	userIDString := r.URL.Query().Get("id")
	userID, err := strconv.Atoi(userIDString)
	if err != nil {
		return nil, err
	}

	user, getUserErr := h.userRepository.GetByID(context.Background(), userID)
	if getUserErr != nil {
		return nil, getUserErr
	}
	userJSON, marshalJSONErr := json.Marshal(user)
	if marshalJSONErr != nil {
		return nil, marshalJSONErr
	}

	return &Response{
		StatusCode: http.StatusOK,
		Body:       userJSON,
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

	userJSON, marshalJSONErr := json.Marshal(users)
	if marshalJSONErr != nil {
		http.Error(w, marshalJSONErr.Error(), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
	_, writeBodyErr := w.Write(userJSON)
	if writeBodyErr != nil {
		http.Error(w, writeBodyErr.Error(), http.StatusInternalServerError)
		return
	}
}
