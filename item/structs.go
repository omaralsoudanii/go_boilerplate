package item

import "mime/multipart"

type File struct {
	Physical multipart.File
	Header   *multipart.FileHeader
}
