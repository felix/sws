package sws

import (
	"crypto/rand"
	"crypto/subtle"
	"encoding/base64"
	"fmt"
	"time"

	"golang.org/x/crypto/argon2"
)

type User struct {
	ID          *int       `json:"id,omitempty"`
	Email       *string    `json:"email,omitempty"`
	FirstName   *string    `json:"first_name,omitempty" db:"first_name"`
	LastName    *string    `json:"last_name,omitempty" db:"last_name"`
	Enabled     bool       `json:"enabled"`
	PwHash      *string    `json:"pw_hash" db:"pw_hash"`
	PwSalt      *string    `json:"pw_salt" db:"pw_salt"`
	LastLoginAt *time.Time `json:"last_login_at" db:"last_login_at"`
	CreatedAt   *time.Time `json:"created_at,omitempty" db:"created_at"`
	UpdatedAt   *time.Time `json:"updated_at,omitempty" db:"updated_at"`
}

const (
	pwMemory  = 64 * 1024
	pwTime    = 1
	pwThreads = 4
	pwLength  = 32
)

func (u *User) SetPassword(pw string) error {
	// Generate a Salt
	saltB := make([]byte, 16)
	if _, err := rand.Read(saltB); err != nil {
		return err
	}

	hashB := generateHash(pw, saltB)

	salt := base64.RawStdEncoding.EncodeToString(saltB)
	hash := base64.RawStdEncoding.EncodeToString(hashB)

	u.PwHash = &hash
	u.PwSalt = &salt
	return nil
}

func (u *User) ValidPassword(test string) error {
	if u.PwHash == nil || u.PwSalt == nil {
		return fmt.Errorf("invalid user")
	}
	var hash1B, hash2B, saltB []byte
	var err error

	if saltB, err = base64.RawStdEncoding.DecodeString(*u.PwSalt); err != nil {
		return err
	}

	hash2B = generateHash(test, saltB)

	if hash1B, err = base64.RawStdEncoding.DecodeString(*u.PwHash); err != nil {
		return err
	}
	ok, err := comparePassword(hash1B, hash2B)
	if err != nil {
		return err
	}
	if !ok {
		return fmt.Errorf("invalid")
	}
	return nil
}

func generateHash(password string, salt []byte) []byte {
	return argon2.IDKey([]byte(password), salt, pwTime, pwMemory, pwThreads, pwLength)
}

func comparePassword(hash1B, hash2B []byte) (bool, error) {
	return (subtle.ConstantTimeCompare(hash1B, hash2B) == 1), nil
}
