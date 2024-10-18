package user

import (
	"context"
	"errors"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgconn"
	"github.com/rs/zerolog/log"
)

var ErrEmailAlreadyExists = errors.New("email already exists")

func findUserByEmail(ctx context.Context, tx pgx.Tx, email string) (User, error) {
	q := `
		SELECT id, email, password, role FROM users WHERE email = $1
	     `
	row := tx.QueryRow(ctx, q, email)

	var user User
	if err := row.Scan(&user.Id, &user.Email, &user.Password, &user.Role); err != nil {
		if err == pgx.ErrNoRows {
			log.Debug().Err(err).Msg("Can't fine any user")
			return User{}, errors.New("User not found")
		}
		return User{}, err
	}
	return user, nil
}

func saveUser(ctx context.Context, tx pgx.Tx, user User) error {
	q := `
	INSERT INTO users (id,  email, password, created_at) VALUES ($1, $2, $3, $4)
	`
	_, err := tx.Exec(ctx, q, user.Id, user.Email, user.Password, user.Created_at)
	if err != nil {
		var pgErr *pgconn.PgError
		if errors.As(err, &pgErr) {
			if pgErr.Code == "23505" && pgErr.ConstraintName == "users_email_key" { // Unique violation on email
				return ErrEmailAlreadyExists
			}
		}
		return err
	}
	return nil
}
