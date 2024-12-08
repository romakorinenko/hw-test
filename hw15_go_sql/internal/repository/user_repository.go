package repository

import (
	"context"
	"errors"
	"fmt"

	"github.com/huandu/go-sqlbuilder"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

type User struct {
	ID       int    `db:"id" fieldtag:"pk" json:"id"`
	Name     string `db:"name" json:"name"`
	Email    string `db:"email" json:"email"`
	Password string `db:"password" json:"password"`
}

type IUserRepository interface {
	Create(ctx context.Context, user *User) (*User, error)
	GetByID(ctx context.Context, userID int) (*User, error)
	GetAll(ctx context.Context) ([]User, error)
	Update(ctx context.Context, user *User) (*User, error)
	DeleteByID(ctx context.Context, userID int) error
}

type UserRepository struct {
	dbPool *pgxpool.Pool
}

func NewUserRepository(dbPool *pgxpool.Pool) *UserRepository {
	return &UserRepository{dbPool: dbPool}
}

const usersTable = "users"

var UserStruct = sqlbuilder.NewStruct(new(User))

func (u *UserRepository) Create(ctx context.Context, user *User) (*User, error) {
	userID, err := u.generateNextUserID(ctx)
	if err != nil {
		return nil, err
	}
	user.ID = userID

	sql, args := UserStruct.InsertInto(usersTable, user).
		BuildWithFlavor(sqlbuilder.PostgreSQL)
	row := u.dbPool.QueryRow(ctx, sql, args...)
	err = row.Scan()
	if err != nil && !errors.Is(err, pgx.ErrNoRows) {
		return nil, err
	}

	return user, nil
}

func (u *UserRepository) GetByID(ctx context.Context, userID int) (*User, error) {
	selectBuilder := UserStruct.SelectFrom(usersTable)
	sql, args := selectBuilder.Where(selectBuilder.Equal("id", userID)).
		BuildWithFlavor(sqlbuilder.PostgreSQL)
	row := u.dbPool.QueryRow(ctx, sql, args...)

	var user User
	err := row.Scan(UserStruct.Addr(&user)...)
	if err != nil {
		return nil, err
	}

	return &user, nil
}

func (u *UserRepository) GetAll(ctx context.Context) ([]User, error) {
	sql, _ := UserStruct.SelectFrom(usersTable).
		OrderBy("id").
		BuildWithFlavor(sqlbuilder.PostgreSQL)

	rows, err := u.dbPool.Query(ctx, sql)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	res := make([]User, 0)
	for rows.Next() {
		var user User
		err := rows.Scan(UserStruct.Addr(&user)...)
		if err != nil {
			return nil, err
		}
		res = append(res, user)
	}

	return res, nil
}

func (u *UserRepository) Update(ctx context.Context, user *User) (*User, error) {
	updateBuilder := sqlbuilder.NewUpdateBuilder()
	sql, args := updateBuilder.Update(usersTable).
		Set(
			updateBuilder.Assign("name", user.Name),
			updateBuilder.Assign("email", user.Email),
			updateBuilder.Assign("password", user.Password),
		).
		Where(updateBuilder.Equal("id", user.ID)).
		BuildWithFlavor(sqlbuilder.PostgreSQL)

	_, err := u.dbPool.Exec(ctx, sql, args...)
	if err != nil {
		return nil, fmt.Errorf("cannot update user: %w", err)
	}

	return user, nil
}

func (u *UserRepository) DeleteByID(ctx context.Context, userID int) error {
	deleteBuilder := UserStruct.DeleteFrom(usersTable)
	sql, args := deleteBuilder.Where(deleteBuilder.Equal("id", userID)).
		BuildWithFlavor(sqlbuilder.PostgreSQL)

	_, err := u.dbPool.Exec(ctx, sql, args...)
	if err != nil {
		return fmt.Errorf("cannot delete user from db: %w", err)
	}

	return nil
}

func (u *UserRepository) generateNextUserID(ctx context.Context) (int, error) {
	rows, err := u.dbPool.Query(ctx, fmt.Sprintf("SELECT nextval('%s')", "users_sequence"))
	if err != nil {
		return 0, err
	}
	defer rows.Close()

	if rows.Next() {
		var id int
		err = rows.Scan(&id)
		if err != nil {
			return 0, err
		}
		return id, nil
	}
	return 0, fmt.Errorf("something was wrong. there is no next user id")
}
