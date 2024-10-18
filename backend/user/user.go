package user

import (
	"time"

	"github.com/alexedwards/argon2id"
	"github.com/oklog/ulid/v2"
	"gopkg.in/guregu/null.v4"
)

type User struct {
	Email      string
	Password   string
	Created_at time.Time
	Updated_at null.Time
	Username   null.String
	Role       null.String
	Id         ulid.ULID
}

func NewUser(email, password string) (*User, error) {
	if err := validateEmail(email); err != nil {
		return &User{}, err
	}
	if err := validatePassword(password); err != nil {
		return &User{}, err
	}
	passwordHash, err := encryptPassword(password)
	if err != nil {
		return &User{}, err
	}

	u := &User{
		Id:         ulid.Make(),
		Email:      email,
		Password:   passwordHash,
		Created_at: time.Now().UTC(),
	}

	return u, nil
}

func (u *User) ChangePassword(newPassword string) error {
	if err := validatePassword(newPassword); err != nil {
		return err
	}
	passwordHash, err := encryptPassword(newPassword)
	if err != nil {
		return err
	}
	u.Password = passwordHash
	u.Updated_at = null.TimeFrom(time.Now())
	return nil
}

func encryptPassword(password string) (hash string, err error) {
	// TODO: Change default params for prod
	hash, err = argon2id.CreateHash(password, argon2id.DefaultParams)
	if err != nil {
		return "", err
	}
	return hash, nil
}

func verifyPassword(userPassword, passwordHash string) (bool, error) {
	// ComparePasswordAndHash performs a constant-time comparison between a
	// plain-text password and Argon2id hash, using the parameters and salt
	// contained in the hash. It returns true if they match, otherwise it returns
	// false.
	match, err := argon2id.ComparePasswordAndHash(userPassword, passwordHash)
	if err != nil {
		return false, err
	}
	return match, nil
}
