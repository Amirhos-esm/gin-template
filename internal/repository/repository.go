package repository

import "template/internal/models"

type DatabaseRepository interface {
	AllUsers() ([]*models.User, error)
	GetUserById(uint64) (*models.User, error)
	GetUserByEmail(string) (*models.User, error)

	Init(databaseURL string) error
	DeInit()
}
