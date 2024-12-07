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
	"github.com/shopspring/decimal"
	"github.com/stretchr/testify/require"
)

func TestOrderHandler(t *testing.T) {
	testDb := test.CreateTestDb(t, "/migrations")
	defer testDb.Close()

	ctx := context.Background()

	host := "http://localhost:8087"
	orderPath := "/orders"
	ordersByIDPath := "/orders/byuserid"
	ordersByEmailPath := "/orders/byuseremail"
	ordersStatsByUserPath := "/orders/user-statistics"
	userPath := "/users"
	productPath := "/products"

	mux := http.NewServeMux()
	orderRepository := repository.NewOrderRepository(testDb.DbPool)
	orderHandler := NewOrderHandler(orderRepository)
	mux.HandleFunc(orderPath, orderHandler.Handle)
	mux.HandleFunc(ordersByIDPath, orderHandler.GetByUserID)
	mux.HandleFunc(ordersByEmailPath, orderHandler.GetByUserEmail)
	mux.HandleFunc(ordersStatsByUserPath, orderHandler.GetStatistics)

	userRepository := repository.NewUserRepository(testDb.DbPool)
	userHandler := NewUserHandler(userRepository)
	mux.HandleFunc(userPath, userHandler.Handle)

	productRepository := repository.NewProductRepository(testDb.DbPool)
	productHandler := NewProductHandler(productRepository)
	mux.HandleFunc(productPath, productHandler.Handle)

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
	postUserResponse, err := client.R().
		SetBody(&repository.User{
			Name:     "User",
			Email:    "User@mail.ru",
			Password: "UserPass",
		}).
		Post(fmt.Sprintf("%s%s", host, userPath))
	require.NoError(t, err)

	var user repository.User
	err = json.Unmarshal(postUserResponse.Body(), &user)
	require.NoError(t, err)

	postFirstProductResponse, err := client.R().
		SetBody(&repository.Product{
			Name:  "apple",
			Price: 25.50,
		}).
		Post(fmt.Sprintf("%s%s", host, productPath))
	require.NoError(t, err)

	var firstProduct repository.Product
	err = json.Unmarshal(postFirstProductResponse.Body(), &firstProduct)
	require.NoError(t, err)

	postSecondProductResponse, err := client.R().
		SetBody(&repository.Product{
			Name:  "cherry",
			Price: 20.00,
		}).
		Post(fmt.Sprintf("%s%s", host, productPath))
	require.NoError(t, err)

	var secondProduct repository.Product
	err = json.Unmarshal(postSecondProductResponse.Body(), &secondProduct)
	require.NoError(t, err)

	postThirdProductResponse, err := client.R().
		SetBody(&repository.Product{
			Name:  "cherry",
			Price: 20.00,
		}).
		Post(fmt.Sprintf("%s%s", host, productPath))
	require.NoError(t, err)

	var thirdProduct repository.Product
	err = json.Unmarshal(postThirdProductResponse.Body(), &thirdProduct)
	require.NoError(t, err)

	postOrderResp, err := client.R().
		SetBody(&repository.Order{
			UserId:      user.ID,
			TotalAmount: float32(decimal.NewFromFloat32(firstProduct.Price).Sub(decimal.NewFromFloat32(secondProduct.Price)).InexactFloat64()),
			ProductIds:  []int{firstProduct.ID, secondProduct.ID},
		}).
		Post(fmt.Sprintf("%s%s", host, orderPath))
	require.NoError(t, err)

	var order repository.Order
	err = json.Unmarshal(postOrderResp.Body(), &order)
	require.NoError(t, err)
	require.Equal(t, 1, order.ID)

	rows, err := testDb.DbPool.Query(ctx, `SELECT COUNT(*) FROM order_products`)
	require.NoError(t, err)

	var count int
	if rows.Next() {
		err = rows.Scan(&count)
	}
	rows.Close()

	require.Equal(t, 2, count)

	getByUserID, err := client.R().
		SetQueryParams(map[string]string{"userId": strconv.Itoa(user.ID)}).
		Get(fmt.Sprintf("%s%s", host, ordersByIDPath))
	require.NoError(t, err)

	var ordersByUserID []repository.Order
	err = json.Unmarshal(getByUserID.Body(), &ordersByUserID)
	require.NoError(t, err)
	require.Len(t, ordersByUserID, 1)
	require.Equal(t, order.ID, ordersByUserID[0].ID)

	getByUserEmail, err := client.R().
		SetQueryParams(map[string]string{"email": user.Email}).
		Get(fmt.Sprintf("%s%s", host, ordersByEmailPath))
	require.NoError(t, err)

	var ordersByUserEmail []repository.Order
	err = json.Unmarshal(getByUserEmail.Body(), &ordersByUserEmail)
	require.NoError(t, err)
	require.Len(t, ordersByUserEmail, 1)
	require.Equal(t, order.ID, ordersByUserEmail[0].ID)

	statistics, err := client.R().
		SetQueryParams(map[string]string{"id": strconv.Itoa(user.ID)}).
		Get(fmt.Sprintf("%s%s", host, ordersStatsByUserPath))
	require.NoError(t, err)

	var statisticsByUser repository.UserStatistics
	err = json.Unmarshal(statistics.Body(), &statisticsByUser)
	require.NoError(t, err)
	require.Equal(t, repository.UserStatistics{
		Name:         user.Name,
		TotalOrders:  1,
		TotalAmount:  45.50,
		AveragePrice: 22.75,
	}, statisticsByUser)

}
