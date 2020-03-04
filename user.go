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
	ID              *int    `json:"id,omitempty"`
	Email           *string `json:"email,omitempty"`
	FirstName       *string `json:"first_name,omitempty" db:"first_name"`
	LastName        *string `json:"last_name,omitempty" db:"last_name"`
	Enabled         bool    `json:"enabled"`
	PwHash          *string `json:"-" db:"pw_hash"`
	PwSalt          *string `json:"-" db:"pw_salt"`
	Password        string  `json:"-" db:"-"`
	PasswordConfirm string  `json:"-" db:"-"`
	Admin           bool
	LastLoginAt     *time.Time `json:"last_login_at" db:"last_login_at"`
	CreatedAt       *time.Time `json:"created_at,omitempty" db:"created_at"`
	UpdatedAt       *time.Time `json:"updated_at,omitempty" db:"updated_at"`
}

const (
	pwMemory  = 64 * 1024
	pwTime    = 1
	pwThreads = 4
	pwLength  = 32
)

func (u *User) Validate() []string {
	var out []string

	if u.FirstName == nil || *u.FirstName == "" {
		out = append(out, fmt.Sprintf("invalid first name"))
	}
	if u.LastName == nil || *u.LastName == "" {
		out = append(out, fmt.Sprintf("invalid last name"))
	}
	if u.Email == nil || *u.Email == "" {
		out = append(out, fmt.Sprintf("invalid email"))
	}
	if u.PwHash == nil || *u.PwHash == "" {
		out = append(out, fmt.Sprint("invalid password"))
	}
	if u.PasswordConfirm != "" {
		if u.Password != u.PasswordConfirm {
			out = append(out, fmt.Sprintf("password confirmation mismatch"))
		} else {
			if err := u.SetPassword(u.Password); err != nil {
				out = append(out, fmt.Sprintf("failed to update password: %s", err))
			}
			u.Password = ""
			u.PasswordConfirm = ""
		}
	}
	return out
}

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
