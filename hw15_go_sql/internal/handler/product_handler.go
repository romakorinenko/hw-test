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

func (h *ProductHandler) create(r *http.Request) (*Response, error) {
	var product repository.Product
	err := json.NewDecoder(r.Body).Decode(&product)
	if err != nil {
		return nil, err
	}

	createdProduct, err := h.productRepository.Create(context.Background(), &product)
	if err != nil {
		return nil, err
	}

	productJSON, err := json.Marshal(createdProduct)
	if err != nil {
		return nil, err
	}
	if err != nil {
		return nil, err
	}
	return &Response{
		StatusCode: http.StatusCreated,
		Body:       productJSON,
	}, nil
}

func (h *ProductHandler) update(r *http.Request) (*Response, error) {
	var product repository.Product
	err := json.NewDecoder(r.Body).Decode(&product)
	if err != nil {
		return nil, err
	}

	updatedProduct, err := h.productRepository.Update(context.Background(), &product)
	if err != nil {
		return nil, err
	}

	productJSON, err := json.Marshal(updatedProduct)
	if err != nil {
		return nil, err
	}
	if err != nil {
		return nil, err
	}
	return &Response{
		StatusCode: http.StatusCreated,
		Body:       productJSON,
	}, nil
}

func (h *ProductHandler) deleteByID(r *http.Request) (*Response, error) {
	productIdString := r.URL.Query().Get("id")
	productId, err := strconv.Atoi(productIdString)
	if err != nil {
		return nil, err
	}

	err = h.productRepository.DeleteByID(context.Background(), productId)
	if err != nil {
		return nil, err
	}

	return &Response{
		StatusCode: http.StatusCreated,
	}, nil
}

func (h *ProductHandler) getByID(r *http.Request) (*Response, error) {
	productIdString := r.URL.Query().Get("id")
	productId, err := strconv.Atoi(productIdString)
	if err != nil {
		return nil, err
	}

	product, err := h.productRepository.GetByID(context.Background(), productId)
	productJSON, err := json.Marshal(product)
	if err != nil {
		return nil, err
	}
	if err != nil {
		return nil, err
	}
	return &Response{
		StatusCode: http.StatusCreated,
		Body:       productJSON,
	}, nil
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
