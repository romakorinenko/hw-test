package handler

import (
	"context"
	"encoding/json"
	"github.com/romakorinenko/hw-test/hw15_go_sql/internal/repository"
	"log"
	"net/http"
	"strconv"
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

	switch r.Method {
	case http.MethodPost:
		h.create(w, r)
	case http.MethodGet:
		h.getByID(w, r)
	case http.MethodPut:
		h.update(w, r)
	case http.MethodDelete:
		h.deleteByID(w, r)
	}
}

func (h *UserHandler) create(w http.ResponseWriter, r *http.Request) {
	var user repository.User
	err := json.NewDecoder(r.Body).Decode(&user)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	createdUser, err := h.userRepository.Create(context.Background(), &user)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	user.ID = createdUser.ID
	userJson, err := json.Marshal(user)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	_, err = w.Write(userJson)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	w.WriteHeader(http.StatusCreated)
}

func (h *UserHandler) update(w http.ResponseWriter, r *http.Request) {
	var user repository.User
	err := json.NewDecoder(r.Body).Decode(&user)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	updatedUser, err := h.userRepository.Update(context.Background(), &user)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	w.WriteHeader(http.StatusCreated)
	userJson, err := json.Marshal(updatedUser)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	_, err = w.Write(userJson)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
}

func (h *UserHandler) deleteByID(w http.ResponseWriter, r *http.Request) {
	userIdString := r.URL.Query().Get("id")
	userId, err := strconv.Atoi(userIdString)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	err = h.userRepository.DeleteById(context.Background(), userId)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
}

func (h *UserHandler) getByID(w http.ResponseWriter, r *http.Request) {
	userIdString := r.URL.Query().Get("id")
	userId, err := strconv.Atoi(userIdString)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	user, err := h.userRepository.GetById(context.Background(), userId)
	userJson, err := json.Marshal(user)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	_, err = w.Write(userJson)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
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
