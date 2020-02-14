package item

import (
	"go_boilerplate/models"
	"mime/multipart"
)

type File struct {
	Physical multipart.File
	Header   *multipart.FileHeader
}

type ItemMapper struct {
	Item *models.Item `db:"item"`
	Category *models.Category `db:"category"`
	User *models.User `db:"user"`
	ItemImages []*models.ItemImages `db:"item_images"`
}

type ItemScanner struct {
	Item *models.Item `db:"item"`
	Category *models.Category `db:"category"`
	User *models.User `db:"user"`
	ItemImages *models.ItemImages `db:"item_images"`
}