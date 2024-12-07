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

func TestUserHandler_CRUD(t *testing.T) {
	testDb := test.CreateTestDb(t, "/migrations")
	defer testDb.Close()

	ctx := context.Background()

	host := "http://localhost:8087"
	userPath := "/users"
	usersPath := "/users/all"

	mux := http.NewServeMux()
	userRepository := repository.NewUserRepository(testDb.DbPool)
	userHandler := NewUserHandler(userRepository)
	mux.HandleFunc(userPath, userHandler.Handle)
	mux.HandleFunc(usersPath, userHandler.GetAll)

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
		SetBody(&repository.User{
			Name:     "User",
			Email:    "User@mail.ru",
			Password: "UserPass",
		}).
		Post(fmt.Sprintf("%s%s", host, userPath))
	require.NoError(t, err)

	createdUser := repository.User{
		ID:       1,
		Name:     "User",
		Email:    "User@mail.ru",
		Password: "UserPass",
	}

	var actualPostUser repository.User
	err = json.Unmarshal(postResp.Body(), &actualPostUser)
	require.NoError(t, err)
	require.Equal(t, createdUser, actualPostUser)

	getResp, err := client.R().
		SetQueryParams(map[string]string{"id": strconv.Itoa(actualPostUser.ID)}).
		Get(fmt.Sprintf("%s%s", host, userPath))
	require.NoError(t, err)

	var actualGetUser repository.User
	err = json.Unmarshal(getResp.Body(), &actualGetUser)
	require.NoError(t, err)
	require.Equal(t, createdUser, actualPostUser)

	putResp, err := client.R().
		SetQueryParams(map[string]string{"id": strconv.Itoa(actualPostUser.ID)}).
		SetBody(&repository.User{
			Name:     actualGetUser.Name,
			Email:    actualGetUser.Email,
			Password: "UserPass11",
		}).
		Put(fmt.Sprintf("%s%s", host, userPath))
	require.NoError(t, err)

	var actualPutUser repository.User
	err = json.Unmarshal(putResp.Body(), &actualPutUser)
	require.NoError(t, err)
	require.Equal(t, "UserPass11", actualPutUser.Password)

	users, err := userRepository.GetAll(ctx)
	require.NoError(t, err)
	require.Equal(t, 1, len(users))

	_, err = client.R().
		SetQueryParams(map[string]string{"id": strconv.Itoa(actualPostUser.ID)}).
		Delete(fmt.Sprintf("%s%s", host, userPath))
	require.NoError(t, err)

	emptyUsers, err := userRepository.GetAll(ctx)
	require.NoError(t, err)
	require.Equal(t, 0, len(emptyUsers))
}
