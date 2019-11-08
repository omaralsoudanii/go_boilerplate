package user

import (
	"context"

	"go_boilerplate/models"
)

// Repository represent the user's repository contract
type Repository interface {
	//Authenticate(ctx context.Context, user *models.User) error
	Insert(ctx context.Context, user *models.User) error
	FetchByName(ctx context.Context, username string) (*models.User, error)
	FetchById(ctx context.Context, ID uint) (*models.User, error)
	StoreSession(ctx context.Context, user *models.User, token string) error
	GetUser(key string) (map[string]string, error)
	DeleteSession(key string) error
}
