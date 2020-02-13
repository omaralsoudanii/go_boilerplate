package user

import (
	"context"

	"go_boilerplate/models"
)

// Repository represent the user's repository contract
type Repository interface {
	//Authenticate(ctx context.Context, user *models.User) error
	Insert(ctx context.Context, user *models.User) (*models.User, error)
	GetByName(ctx context.Context, username string) (*models.User, error)
	GetByID(ctx context.Context, id uint) (*models.User, error)
	StoreSession(user *models.User, key string, token string) error
	GetUser(key string) (map[string]string, error)
	DeleteSession(key string) error
}
