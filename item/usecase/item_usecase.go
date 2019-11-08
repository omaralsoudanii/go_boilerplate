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

	_lib "go_boilerplate/lib"
)

var log = _lib.GetLogger()

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

func (i *itemUseCase) Fetch(c context.Context, num int64) ([]*models.Item, error) {

	ctx, cancel := context.WithTimeout(c, i.contextTimeout)
	defer cancel()
	listItems, err := i.itemRepo.Fetch(ctx, num)
	if err != nil {
		return nil, err
	}

	return listItems, nil
}

func (i *itemUseCase) GetByID(c context.Context, id int64) (*models.Item, error) {

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

func (i *itemUseCase) GetByTitle(c context.Context, name string) (*models.Item, error) {

	ctx, cancel := context.WithTimeout(c, i.contextTimeout)
	defer cancel()
	res, err := i.itemRepo.GetByTitle(ctx, name)
	if err != nil {
		return nil, err
	}

	return res, nil
}

func (i *itemUseCase) Store(c context.Context, item *models.Item, images []item.File) (uint, error) {

	ctx, cancel := context.WithTimeout(c, i.contextTimeout)
	defer cancel()
	userContext, ok := c.Value("user").(*user.ContextData)
	user, err := i.userRepo.GetUser(userContext.UserName)
	if err != nil {
		return 0, err
	}
	UserID, err := strconv.Atoi(user["id"])
	item.UserID = UserID
	if !ok {
		return 0, errors.New("context_retrieve_user_err")
	}
	var fileNames []string
	for _, image := range images {
		fileName, err := lib.UploadFile(ctx, image, "items")
		fileNames = append(fileNames, fileName)
		if err != nil {
			log.Error(err)
			return 0, err
		}
	}
	id, err := i.itemRepo.Store(ctx, item, fileNames)
	if err != nil {
		return 0, err
	}
	return id, nil
}

func (i *itemUseCase) Delete(c context.Context, id int64) error {
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
