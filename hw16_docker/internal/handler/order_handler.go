package handler

import (
	"context"
	"encoding/json"
	"log"
	"net/http"
	"strconv"

	"github.com/romakorinenko/hw-test/hw16_docker/internal/repository"
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

func NewOrderHandler(orderRepository repository.IOrderRepository) *OrderHandler {
	return &OrderHandler{orderRepository: orderRepository}
}

func (h *OrderHandler) Handle(w http.ResponseWriter, r *http.Request) {
	log.Println("received request: ", r.Method, r.URL.Path)

	var response *Response
	var handlingErr error

	switch r.Method {
	case http.MethodPut:
		http.Error(w, "", http.StatusMethodNotAllowed)
		return
	case http.MethodPost:
		response, handlingErr = h.create(r)
	case http.MethodDelete:
		response, handlingErr = h.deleteByID(r)
	case http.MethodGet:
		http.Error(w, "", http.StatusMethodNotAllowed)
		return
	}

	if handlingErr != nil {
		http.Error(w, handlingErr.Error(), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(response.StatusCode)
	if _, err := w.Write(response.Body); handlingErr != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}

func (h *OrderHandler) create(r *http.Request) (*Response, error) {
	var order repository.Order
	err := json.NewDecoder(r.Body).Decode(&order)
	if err != nil {
		return nil, err
	}

	createdOrder, creationOrderErr := h.orderRepository.Create(context.Background(), &order)
	if creationOrderErr != nil {
		return nil, creationOrderErr
	}

	orderJSON, marshalJSONErr := json.Marshal(createdOrder)
	if marshalJSONErr != nil {
		return nil, marshalJSONErr
	}

	return &Response{
		StatusCode: http.StatusCreated,
		Body:       orderJSON,
	}, nil
}

func (h *OrderHandler) deleteByID(r *http.Request) (*Response, error) {
	orderIDString := r.URL.Query().Get("id")
	productID, err := strconv.Atoi(orderIDString)
	if err != nil {
		return nil, err
	}

	deleteOrderErr := h.orderRepository.DeleteByID(context.Background(), productID)
	if deleteOrderErr != nil {
		return nil, deleteOrderErr
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

	userIDString := r.URL.Query().Get("userId")
	userID, err := strconv.Atoi(userIDString)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	orders, getUserErr := h.orderRepository.GetByUserID(context.Background(), userID)
	if getUserErr != nil {
		http.Error(w, getUserErr.Error(), http.StatusInternalServerError)
		return
	}

	ordersJSON, unmarshalJSONErr := json.Marshal(orders)
	if unmarshalJSONErr != nil {
		http.Error(w, unmarshalJSONErr.Error(), http.StatusInternalServerError)
		return
	}
	_, writeBodyErr := w.Write(ordersJSON)
	if err != nil {
		http.Error(w, writeBodyErr.Error(), http.StatusInternalServerError)
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

	ordersJSON, marshalJSONErr := json.Marshal(orders)
	if marshalJSONErr != nil {
		http.Error(w, marshalJSONErr.Error(), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
	_, writeJSONErr := w.Write(ordersJSON)
	if writeJSONErr != nil {
		http.Error(w, writeJSONErr.Error(), http.StatusInternalServerError)
		return
	}
}

func (h *OrderHandler) GetStatistics(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, "", http.StatusMethodNotAllowed)
		return
	}

	userIDString := r.URL.Query().Get("id")
	userID, err := strconv.Atoi(userIDString)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	userStat, getStatErr := h.orderRepository.GetStatisticsByID(context.Background(), userID)
	if getStatErr != nil {
		http.Error(w, getStatErr.Error(), http.StatusInternalServerError)
		return
	}

	userStatJSON, marshalJSONErr := json.Marshal(userStat)
	if marshalJSONErr != nil {
		http.Error(w, marshalJSONErr.Error(), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
	_, writeBodyErr := w.Write(userStatJSON)
	if writeBodyErr != nil {
		http.Error(w, writeBodyErr.Error(), http.StatusInternalServerError)
		return
	}
}
