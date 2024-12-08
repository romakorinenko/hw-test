package handler

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"
	"testing"
	"time"

	"github.com/go-resty/resty/v2"
	"github.com/romakorinenko/hw-test/hw15_go_sql/internal/repository"
	"github.com/romakorinenko/hw-test/hw15_go_sql/test"
	"github.com/stretchr/testify/require"
)

const productURLHost = "http://localhost:8087"

func TestProductHandler_CRUD(t *testing.T) {
	testDB := test.CreateDBForTest(t, "/migrations")
	defer testDB.Close()

	ctx := context.Background()

	productPath := "/products"
	productsPath := "/products/all"

	mux := http.NewServeMux()
	productRepository := repository.NewProductRepository(testDB.DBPool)
	productHandler := NewProductHandler(productRepository)
	mux.HandleFunc(productPath, productHandler.Handle)
	mux.HandleFunc(productsPath, productHandler.GetAll)

	server := &http.Server{
		Addr:              ":8087",
		Handler:           mux,
		ReadHeaderTimeout: 5 * time.Second,
	}
	go func() {
		_ = server.ListenAndServe()
	}()
	defer func() {
		_ = server.Shutdown(ctx)
	}()

	client := resty.New()
	postResp, err := client.R().
		SetBody(&repository.Product{
			Name:  "apple",
			Price: 25.50,
		}).
		Post(fmt.Sprintf("%s%s", productURLHost, productPath))
	require.NoError(t, err)

	createdProduct := repository.Product{
		ID:    1,
		Name:  "apple",
		Price: 25.50,
	}

	var actualPostProduct repository.Product
	err = json.Unmarshal(postResp.Body(), &actualPostProduct)
	require.NoError(t, err)
	require.Equal(t, createdProduct, actualPostProduct)

	getResp, err := client.R().
		SetQueryParams(map[string]string{"id": strconv.Itoa(actualPostProduct.ID)}).
		Get(fmt.Sprintf("%s%s", productURLHost, productPath))
	require.NoError(t, err)

	var actualGetProduct repository.Product
	err = json.Unmarshal(getResp.Body(), &actualGetProduct)
	require.NoError(t, err)
	require.Equal(t, createdProduct, actualPostProduct)

	putResp, err := client.R().
		SetQueryParams(map[string]string{"id": strconv.Itoa(actualPostProduct.ID)}).
		SetBody(&repository.Product{
			Name:  actualGetProduct.Name,
			Price: 20.50,
		}).
		Put(fmt.Sprintf("%s%s", productURLHost, productPath))
	require.NoError(t, err)

	var actualPutProduct repository.Product
	err = json.Unmarshal(putResp.Body(), &actualPutProduct)
	require.NoError(t, err)
	require.Equal(t, float32(20.50), actualPutProduct.Price)

	products, err := productRepository.GetAll(ctx)
	require.NoError(t, err)
	require.Equal(t, 1, len(products))

	_, err = client.R().
		SetQueryParams(map[string]string{"id": strconv.Itoa(actualPostProduct.ID)}).
		Delete(fmt.Sprintf("%s%s", productURLHost, productPath))
	require.NoError(t, err)

	emptyProducts, err := productRepository.GetAll(ctx)
	require.NoError(t, err)
	require.Equal(t, 0, len(emptyProducts))
}
