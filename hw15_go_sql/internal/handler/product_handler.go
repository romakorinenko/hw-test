package handler

import (
	"context"
	"encoding/json"
	"log"
	"net/http"
	"strconv"

	"github.com/romakorinenko/hw-test/hw15_go_sql/internal/repository"
)

type IProductHandler interface {
	Handle(w http.ResponseWriter, r *http.Request)
	GetAll(w http.ResponseWriter, r *http.Request)
}

type ProductHandler struct {
	productRepository repository.IProductRepository
}

func NewProductHandler(productRepository repository.IProductRepository) IProductHandler {
	return &ProductHandler{productRepository: productRepository}
}

func (h *ProductHandler) Handle(w http.ResponseWriter, r *http.Request) {
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

func (h *ProductHandler) create(w http.ResponseWriter, r *http.Request) {
	var product repository.Product
	err := json.NewDecoder(r.Body).Decode(&product)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	createdProduct, err := h.productRepository.Create(context.Background(), &product)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	productJson, err := json.Marshal(createdProduct)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	_, err = w.Write(productJson)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	w.WriteHeader(http.StatusCreated)
}

func (h *ProductHandler) update(w http.ResponseWriter, r *http.Request) {
	var product repository.Product
	err := json.NewDecoder(r.Body).Decode(&product)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	updatedProduct, err := h.productRepository.Update(context.Background(), &product)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	w.WriteHeader(http.StatusCreated)
	productJson, err := json.Marshal(updatedProduct)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	_, err = w.Write(productJson)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
}

func (h *ProductHandler) deleteByID(w http.ResponseWriter, r *http.Request) {
	productIdString := r.URL.Query().Get("id")
	productId, err := strconv.Atoi(productIdString)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	err = h.productRepository.DeleteByID(context.Background(), productId)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
}

func (h *ProductHandler) getByID(w http.ResponseWriter, r *http.Request) {
	productIdString := r.URL.Query().Get("id")
	productId, err := strconv.Atoi(productIdString)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	product, err := h.productRepository.GetByID(context.Background(), productId)
	productJson, err := json.Marshal(product)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	_, err = w.Write(productJson)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
}

func (h *ProductHandler) GetAll(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, "", http.StatusMethodNotAllowed)
		return
	}

	products, err := h.productRepository.GetAll(context.Background())
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	productJson, err := json.Marshal(products)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
	_, err = w.Write(productJson)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
}
