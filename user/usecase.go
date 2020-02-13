package user

import (
	"context"

	"go_boilerplate/models"
)

// UseCase represent the job use cases
type UseCase interface {
	Register(ctx context.Context, user *models.User) (*models.User,string,string, error)
	SignIn(ctx context.Context, user *models.User) (string, string, error)
	Refresh(ctx context.Context, refreshToken string) (string, error)
	SignOut(ctx context.Context) error
}
