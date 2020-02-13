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
	Item models.Item
	Category models.Category
	User models.User
	ItemImages []models.ItemImages
}