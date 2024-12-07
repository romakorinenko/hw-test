package handler

import (
	"context"
	"encoding/json"
	"log"
	"net/http"
	"strconv"

	"github.com/romakorinenko/hw-test/hw15_go_sql/internal/repository"
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

	var response *Response
	var err error

	switch r.Method {
	case http.MethodPost:
		response, err = h.create(r)
	case http.MethodGet:
		http.Error(w, "", http.StatusMethodNotAllowed)
		return
	case http.MethodPut:
		http.Error(w, "", http.StatusMethodNotAllowed)
		return
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

func (h *OrderHandler) create(r *http.Request) (*Response, error) {
	var order repository.Order
	err := json.NewDecoder(r.Body).Decode(&order)
	if err != nil {
		return nil, err
	}

	createdOrder, err := h.orderRepository.Create(context.Background(), &order)
	if err != nil {
		return nil, err
	}

	orderJSON, err := json.Marshal(createdOrder)
	if err != nil {
		return nil, err
	}

	return &Response{
		StatusCode: http.StatusCreated,
		Body:       orderJSON,
	}, nil
}

func (h *OrderHandler) deleteByID(r *http.Request) (*Response, error) {
	orderIdString := r.URL.Query().Get("id")
	productId, err := strconv.Atoi(orderIdString)
	if err != nil {
		return nil, err
	}

	err = h.orderRepository.DeleteById(context.Background(), productId)
	if err != nil {
		return nil, err
	}

	return &Response{
		StatusCode: http.StatusCreated,
	}, nil
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
