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

const orderURLHost = "http://localhost:8087"

func TestOrderHandler(t *testing.T) {
	testDB := test.CreateDBForTest(t, "/migrations")
	defer testDB.Close()

	ctx := context.Background()

	orderPath := "/orders"
	ordersByIDPath := "/orders/byuserid"
	ordersByEmailPath := "/orders/byuseremail"
	ordersStatsByUserPath := "/orders/user-statistics"
	userPath := "/users"
	productPath := "/products"

	mux := http.NewServeMux()
	orderRepository := repository.NewOrderRepository(testDB.DBPool)
	orderHandler := NewOrderHandler(orderRepository)
	mux.HandleFunc(orderPath, orderHandler.Handle)
	mux.HandleFunc(ordersByIDPath, orderHandler.GetByUserID)
	mux.HandleFunc(ordersByEmailPath, orderHandler.GetByUserEmail)
	mux.HandleFunc(ordersStatsByUserPath, orderHandler.GetStatistics)

	userRepository := repository.NewUserRepository(testDB.DBPool)
	userHandler := NewUserHandler(userRepository)
	mux.HandleFunc(userPath, userHandler.Handle)

	productRepository := repository.NewProductRepository(testDB.DBPool)
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
	postUserResponse, postUserResponseErr := client.R().
		SetBody(&repository.User{
			Name:     "User",
			Email:    "User@mail.ru",
			Password: "UserPass",
		}).
		Post(fmt.Sprintf("%s%s", orderURLHost, userPath))
	require.NoError(t, postUserResponseErr)

	var user repository.User
	userErr := json.Unmarshal(postUserResponse.Body(), &user)
	require.NoError(t, userErr)

	postFirstProductResponse, postFirstProductResponseErr := client.R().
		SetBody(&repository.Product{
			Name:  "apple",
			Price: 25.50,
		}).
		Post(fmt.Sprintf("%s%s", orderURLHost, productPath))
	require.NoError(t, postFirstProductResponseErr)

	var firstProduct repository.Product
	firstProductErr := json.Unmarshal(postFirstProductResponse.Body(), &firstProduct)
	require.NoError(t, firstProductErr)

	postSecondProductResponse, postSecondProductResponseErr := client.R().
		SetBody(&repository.Product{
			Name:  "cherry",
			Price: 20.00,
		}).
		Post(fmt.Sprintf("%s%s", orderURLHost, productPath))
	require.NoError(t, postSecondProductResponseErr)

	var secondProduct repository.Product
	secondProductUnmarshalErr := json.Unmarshal(postSecondProductResponse.Body(), &secondProduct)
	require.NoError(t, secondProductUnmarshalErr)

	postOrderResp, postOrderRespErr := client.R().
		SetBody(&repository.Order{
			UserID:      user.ID,
			TotalAmount: float32(decimal.NewFromFloat32(firstProduct.Price).Sub(decimal.NewFromFloat32(secondProduct.Price)).InexactFloat64()),
			ProductIDs:  []int{firstProduct.ID, secondProduct.ID},
		}).
		Post(fmt.Sprintf("%s%s", orderURLHost, orderPath))
	require.NoError(t, postOrderRespErr)

	var order repository.Order
	orderUnmarshalErr := json.Unmarshal(postOrderResp.Body(), &order)
	require.NoError(t, orderUnmarshalErr)
	require.Equal(t, 1, order.ID)

	rows, rowsErr := testDB.DBPool.Query(ctx, `SELECT COUNT(*) FROM order_products`)
	require.NoError(t, rowsErr)

	var count int
	if rows.Next() {
		rowScanErr := rows.Scan(&count)
		require.NoError(t, rowScanErr)
	}
	rows.Close()

	require.Equal(t, 2, count)

	getByUserID, getByUserIDErr := client.R().
		SetQueryParams(map[string]string{"userId": strconv.Itoa(user.ID)}).
		Get(fmt.Sprintf("%s%s", orderURLHost, ordersByIDPath))
	require.NoError(t, getByUserIDErr)

	var ordersByUserID []repository.Order
	ordersByUserIDErr := json.Unmarshal(getByUserID.Body(), &ordersByUserID)
	require.NoError(t, ordersByUserIDErr)
	require.Len(t, ordersByUserID, 1)
	require.Equal(t, order.ID, ordersByUserID[0].ID)

	getByUserEmail, getByUserEmailErr := client.R().
		SetQueryParams(map[string]string{"email": user.Email}).
		Get(fmt.Sprintf("%s%s", orderURLHost, ordersByEmailPath))
	require.NoError(t, getByUserEmailErr)

	var ordersByUserEmail []repository.Order
	ordersByUserEmailUnmarshalErr := json.Unmarshal(getByUserEmail.Body(), &ordersByUserEmail)
	require.NoError(t, ordersByUserEmailUnmarshalErr)
	require.Len(t, ordersByUserEmail, 1)
	require.Equal(t, order.ID, ordersByUserEmail[0].ID)

	statistics, statisticsErr := client.R().
		SetQueryParams(map[string]string{"id": strconv.Itoa(user.ID)}).
		Get(fmt.Sprintf("%s%s", orderURLHost, ordersStatsByUserPath))
	require.NoError(t, statisticsErr)

	var statisticsByUser repository.UserStatistics
	statisticsByUserUnmarshalErr := json.Unmarshal(statistics.Body(), &statisticsByUser)
	require.NoError(t, statisticsByUserUnmarshalErr)
	require.Equal(t, repository.UserStatistics{
		Name:         user.Name,
		TotalOrders:  1,
		TotalAmount:  45.50,
		AveragePrice: 22.75,
	}, statisticsByUser)
}
