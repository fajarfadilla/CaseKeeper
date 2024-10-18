package user

import (
	"context"

	"github.com/oklog/ulid/v2"
)

func createUser(ctx context.Context, email, password string) (id ulid.ULID, err error) {
	user, err := NewUser(email, password)
	if err != nil {
		return
	}
	tx, err := pool.Begin(ctx)
	if err != nil {
		return
	}
	err = saveUser(ctx, tx, *user)
	if err != nil {
		tx.Rollback(ctx)
		return
	}
	err = tx.Commit(ctx)
	if err != nil {
		return
	}
	return user.Id, nil
}

func findUser(ctx context.Context, email string) (user User, err error) {
	tx, err := pool.Begin(ctx)
	if err != nil {
		return
	}
	user, err = findUserByEmail(ctx, tx, email)
	if err != nil {
		return User{}, err
	}
	err = tx.Commit(ctx)
	return
}
