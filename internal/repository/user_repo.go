package repository

import (
	"banking-api/internal/db"
	"banking-api/internal/model"
	"context"
)

type UserRepository struct{}

func (r *UserRepository) Create(ctx context.Context, user *model.User) error {
	query := `
		INSERT INTO USERS (email, username, password)
		VALUES ($1, $2, $3)
		RETURNING id;
	`
	return db.DB.QueryRowContext(
		ctx, query, user.Email, user.Username, user.Password).Scan(&user.ID)
}

func (r *UserRepository) ExistsByEmailOrUsername(ctx context.Context, email, username string) (bool, error) {
	query := `SELECT COUNT (*) FROM users WHERE email=$1 OR username=$2`
	var count int
	err := db.DB.QueryRowContext(ctx, query, email, username).Scan(&count)
	return count > 0, err
}

func (r *UserRepository) FindByEmail(ctx context.Context, email string) (*model.User, error) {
	query := `
		SELECT id, email, username, password
		FROM users
		WHERE email=$1
	`
	row := db.DB.QueryRowContext(ctx, query, email)

	var user model.User
	err := row.Scan(&user.ID, &user.Email, &user.Username, &user.Password)
	if err != nil {
		return nil, err
	}
	return &user, nil
}
