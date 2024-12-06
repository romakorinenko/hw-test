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
	GetByUserID(w http.ResponseWriter, r *http.Request)
	GetByUserEmail(w http.ResponseWriter, r *http.Request)
	GetStatistics(w http.ResponseWriter, r *http.Request)
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

func (h *OrderHandler) GetByUserID(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, "", http.StatusInternalServerError)
		return
	}

	userIdString := r.URL.Query().Get("userId")
	userId, err := strconv.Atoi(userIdString)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	orders, err := h.orderRepository.GetByUserID(context.Background(), userId)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	ordersJson, err := json.Marshal(orders)
	_, err = w.Write(ordersJson)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	w.WriteHeader(http.StatusOK)
}

func (h *OrderHandler) GetByUserEmail(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, "", http.StatusMethodNotAllowed)
		return
	}

	userEmail := r.URL.Query().Get("email")

	orders, err := h.orderRepository.GetByUserEmail(context.Background(), userEmail)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	ordersJson, err := json.Marshal(orders)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
	_, err = w.Write(ordersJson)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
}

func (h *OrderHandler) GetStatistics(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, "", http.StatusMethodNotAllowed)
		return
	}

	userIdString := r.URL.Query().Get("id")
	userId, err := strconv.Atoi(userIdString)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	userStat, err := h.orderRepository.GetStatisticsByID(context.Background(), userId)

	userStatJson, err := json.Marshal(userStat)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
	_, err = w.Write(userStatJson)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
}
