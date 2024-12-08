//nolint:dupl
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

func NewProductHandler(productRepository repository.IProductRepository) *ProductHandler {
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
	if _, writeBodyErr := w.Write(response.Body); err != nil {
		http.Error(w, writeBodyErr.Error(), http.StatusInternalServerError)
		return
	}
}

func (h *ProductHandler) create(r *http.Request) (*Response, error) {
	var product repository.Product
	err := json.NewDecoder(r.Body).Decode(&product)
	if err != nil {
		return nil, err
	}

	createdProduct, createProductErr := h.productRepository.Create(context.Background(), &product)
	if createProductErr != nil {
		return nil, createProductErr
	}

	productJSON, marshalJSONErr := json.Marshal(createdProduct)
	if marshalJSONErr != nil {
		return nil, marshalJSONErr
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

	updatedProduct, updateProductErr := h.productRepository.Update(context.Background(), &product)
	if updateProductErr != nil {
		return nil, updateProductErr
	}

	productJSON, marshalJSONErr := json.Marshal(updatedProduct)
	if marshalJSONErr != nil {
		return nil, marshalJSONErr
	}
	return &Response{
		StatusCode: http.StatusCreated,
		Body:       productJSON,
	}, nil
}

func (h *ProductHandler) deleteByID(r *http.Request) (*Response, error) {
	productIDString := r.URL.Query().Get("id")
	productID, err := strconv.Atoi(productIDString)
	if err != nil {
		return nil, err
	}

	deleteProductErr := h.productRepository.DeleteByID(context.Background(), productID)
	if deleteProductErr != nil {
		return nil, deleteProductErr
	}

	return &Response{
		StatusCode: http.StatusCreated,
	}, nil
}

func (h *ProductHandler) getByID(r *http.Request) (*Response, error) {
	productIDString := r.URL.Query().Get("id")
	productID, err := strconv.Atoi(productIDString)
	if err != nil {
		return nil, err
	}

	product, getProductErr := h.productRepository.GetByID(context.Background(), productID)
	if getProductErr != nil {
		return nil, getProductErr
	}
	productJSON, marshalJSONErr := json.Marshal(product)
	if marshalJSONErr != nil {
		return nil, marshalJSONErr
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

	productJSON, marshalJSONErr := json.Marshal(products)
	if marshalJSONErr != nil {
		http.Error(w, marshalJSONErr.Error(), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
	_, writeBodyErr := w.Write(productJSON)
	if writeBodyErr != nil {
		http.Error(w, writeBodyErr.Error(), http.StatusInternalServerError)
		return
	}
}
