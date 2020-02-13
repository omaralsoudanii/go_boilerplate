package item

import (
	"context"

	"go_boilerplate/models"
)

// UseCase represent the item usecases
type UseCase interface {
	GetAll(ctx context.Context, num uint) ([]ItemMapper, error)
	GetByID(ctx context.Context, id uint) (*ItemMapper, error)
	GetByTitle(ctx context.Context, title string) (*ItemMapper, error)
	Update(ctx context.Context, item *models.Item) error
	Store(ctx context.Context, item *models.Item, images []File) (int64, error)
	Delete(ctx context.Context, id uint) error
}
