package models

import (
	"database/sql/driver"
	"encoding/json"
	"errors"
	"fmt"

	"github.com/google/uuid"
	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
)

type JSONB map[string]interface{}

type User struct {
	gorm.Model
	Name     string    `json:"name" gorm:"size:255;not null"`
	LastName string    `json:"last_name" gorm:"size:255;not null"`
	UUID     uuid.UUID `json:"uuid" gorm:"type:char(36);not null;unique"`
	Password string    `json:"-" gorm:"size:255;not null"`
	Email    string    `json:"email" gorm:"size:255;not null;unique"`
	MetaData JSONB     `json:"meta_data" gorm:"type:json"` // Custom serialization required
}

// Implement driver.Valuer for JSONB (serialization to DB)
func (j JSONB) Value() (driver.Value, error) {
	return json.Marshal(j)
}

// Implement sql.Scanner for JSONB (deserialization from DB)
func (j *JSONB) Scan(value interface{}) error {
	if value == nil {
		*j = make(JSONB)
		return nil
	}
	data, ok := value.([]byte)
	if !ok {
		return fmt.Errorf("failed to scan JSONB: %v", value)
	}
	return json.Unmarshal(data, j)
}

func (u *User) PasswordMatches(plainText string) (bool, error) {
	err := bcrypt.CompareHashAndPassword([]byte(u.Password), []byte(plainText))
	if err != nil {
		switch {
		case errors.Is(err, bcrypt.ErrMismatchedHashAndPassword):
			// invalid password
			return false, nil
		default:
			return false, err
		}
	}

	return true, nil
}
