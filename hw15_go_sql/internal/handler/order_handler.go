package handler

import (
	"context"
	"encoding/json"
	"github.com/romakorinenko/hw-test/hw15_go_sql/internal/repository"
	"log"
	"net/http"
	"strconv"
)

type IOrderHandler interface {
	Handle(w http.ResponseWriter, r *http.Request)
}

type OrderHandler struct {
	orderRepository repository.IOrderRepository
}

func NewOrderHandler(orderRepository repository.IOrderRepository) IOrderHandler {
	return &OrderHandler{orderRepository: orderRepository}
}

func (h *OrderHandler) Handle(w http.ResponseWriter, r *http.Request) {
	log.Println("received request: ", r.Method, r.URL.Path)

	switch r.Method {
	case http.MethodPost:
		h.create(w, r)
	case http.MethodGet:
		http.Error(w, "", http.StatusMethodNotAllowed)
		return
	case http.MethodPut:
		http.Error(w, "", http.StatusMethodNotAllowed)
		return
	case http.MethodDelete:
		h.deleteByID(w, r)
	}
}

func (h *OrderHandler) create(w http.ResponseWriter, r *http.Request) {
	var order repository.Order
	err := json.NewDecoder(r.Body).Decode(&order)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	createdOrder, err := h.orderRepository.Create(context.Background(), &order)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	orderJson, err := json.Marshal(createdOrder)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	_, err = w.Write(orderJson)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	w.WriteHeader(http.StatusCreated)
}

func (h *OrderHandler) deleteByID(w http.ResponseWriter, r *http.Request) {
	orderIdString := r.URL.Query().Get("id")
	productId, err := strconv.Atoi(orderIdString)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	err = h.orderRepository.DeleteById(context.Background(), productId)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
}
