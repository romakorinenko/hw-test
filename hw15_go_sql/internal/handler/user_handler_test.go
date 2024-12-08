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

const userURLHost = "http://localhost:8087"

func TestUserHandler_CRUD(t *testing.T) {
	testDB := test.CreateDBForTest(t, "/migrations")
	defer testDB.Close()

	ctx := context.Background()

	userPath := "/users"
	usersPath := "/users/all"

	mux := http.NewServeMux()
	userRepository := repository.NewUserRepository(testDB.DBPool)
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
	postResp, postRespErr := client.R().
		SetBody(&repository.User{
			Name:     "User",
			Email:    "User@mail.ru",
			Password: "UserPass",
		}).
		Post(fmt.Sprintf("%s%s", userURLHost, userPath))
	require.NoError(t, postRespErr)

	createdUser := repository.User{
		ID:       1,
		Name:     "User",
		Email:    "User@mail.ru",
		Password: "UserPass",
	}

	var actualPostUser repository.User
	actualPostUserUnmarshalErr := json.Unmarshal(postResp.Body(), &actualPostUser)
	require.NoError(t, actualPostUserUnmarshalErr)
	require.Equal(t, createdUser, actualPostUser)

	getResp, getRespErr := client.R().
		SetQueryParams(map[string]string{"id": strconv.Itoa(actualPostUser.ID)}).
		Get(fmt.Sprintf("%s%s", userURLHost, userPath))
	require.NoError(t, getRespErr)

	var actualGetUser repository.User
	actualGetUserErr := json.Unmarshal(getResp.Body(), &actualGetUser)
	require.NoError(t, actualGetUserErr)
	require.Equal(t, createdUser, actualPostUser)

	putResp, putRespErr := client.R().
		SetQueryParams(map[string]string{"id": strconv.Itoa(actualPostUser.ID)}).
		SetBody(&repository.User{
			Name:     actualGetUser.Name,
			Email:    actualGetUser.Email,
			Password: "UserPass11",
		}).
		Put(fmt.Sprintf("%s%s", userURLHost, userPath))
	require.NoError(t, putRespErr)

	var actualPutUser repository.User
	actualPutUserUnmarshalErr := json.Unmarshal(putResp.Body(), &actualPutUser)
	require.NoError(t, actualPutUserUnmarshalErr)
	require.Equal(t, "UserPass11", actualPutUser.Password)

	users, err := userRepository.GetAll(ctx)
	require.NoError(t, err)
	require.Equal(t, 1, len(users))

	_, err = client.R().
		SetQueryParams(map[string]string{"id": strconv.Itoa(actualPostUser.ID)}).
		Delete(fmt.Sprintf("%s%s", userURLHost, userPath))
	require.NoError(t, err)

	emptyUsers, getUsersErr := userRepository.GetAll(ctx)
	require.NoError(t, getUsersErr)
	require.Equal(t, 0, len(emptyUsers))
}
