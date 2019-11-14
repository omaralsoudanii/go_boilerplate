package usecase

import (
	"context"
	"errors"
	"go_boilerplate/item"
	"go_boilerplate/lib"
	"go_boilerplate/models"
	"go_boilerplate/user"
	"strconv"
	"time"
)

type itemUseCase struct {
	itemRepo       item.Repository
	contextTimeout time.Duration
	userRepo       user.Repository
}

// NewItemUseCase will create new an itemUseCase object representation of item.UseCase interface
func NewItemUseCase(i item.Repository, u user.Repository, timeout time.Duration) item.UseCase {
	return &itemUseCase{
		itemRepo:       i,
		userRepo:       u,
		contextTimeout: timeout,
	}
}

func (i *itemUseCase) GetAll(c context.Context, num uint) ([]*models.Item, error) {

	ctx, cancel := context.WithTimeout(c, i.contextTimeout)
	defer cancel()
	listItems, err := i.itemRepo.GetAll(ctx, num)
	if err != nil {
		return nil, err
	}

	return listItems, nil
}

func (i *itemUseCase) GetByID(c context.Context, id uint) (*models.Item, error) {

	ctx, cancel := context.WithTimeout(c, i.contextTimeout)
	defer cancel()

	res, err := i.itemRepo.GetByID(ctx, id)
	if err != nil {
		return nil, err
	}
	return res, nil
}

func (i *itemUseCase) Update(c context.Context, item *models.Item) error {

	ctx, cancel := context.WithTimeout(c, i.contextTimeout)
	defer cancel()

	return i.itemRepo.Update(ctx, item)
}

func (i *itemUseCase) GetByTitle(c context.Context, GetByTitle string) (*models.Item, error) {

	ctx, cancel := context.WithTimeout(c, i.contextTimeout)
	defer cancel()
	res, err := i.itemRepo.GetByTitle(ctx, GetByTitle)
	if err != nil {
		return nil, err
	}

	return res, nil
}

func (i *itemUseCase) Store(c context.Context, item *models.Item, images []item.File) (uint, error) {

	ctx, cancel := context.WithTimeout(c, i.contextTimeout)
	defer cancel()
	userContext, ok := c.Value("user").(*user.ContextData)
	if !ok {
		return 0, errors.New("context_retrieve_user_err")
	}
	user, err := i.userRepo.GetUser(userContext.UserName)
	if err != nil {
		return 0, err
	}

	UserID, err := strconv.Atoi(user["id"])
	if err != nil {
		return 0, err
	}
	item.UserID = UserID

	var fileNames []string
	for _, image := range images {
		fileName, err := lib.UploadFile(ctx, image, "items")
		if err != nil {
			return 0, err
		}
		fileNames = append(fileNames, fileName)
	}
	id, err := i.itemRepo.Store(ctx, item, fileNames)
	if err != nil {
		return 0, err
	}
	return id, nil
}

func (i *itemUseCase) Delete(c context.Context, id uint) error {
	ctx, cancel := context.WithTimeout(c, i.contextTimeout)
	defer cancel()
	//jobExist, err := j.jobRepo.GetByID(ctx, id)
	//if err != nil {
	//	return err
	//}
	//if jobExist == nil {
	//	return models.ErrNotFound
	//}
	return i.itemRepo.Delete(ctx, id)
}
