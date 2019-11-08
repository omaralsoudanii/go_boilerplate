package lib

import (
	"context"
	"crypto/sha256"
	"encoding/hex"
	"errors"
	"fmt"
	"go_boilerplate/item"
	"go_boilerplate/user"
	"io"
	"os"
	"strings"

	uuid "github.com/satori/go.uuid"
)

var allowedContentTypes = []string{
	"image/png",
	"image/jpeg",
	"image/jpg",
	"image/gif"}

// UploadFile ... uploads a file
func UploadFile(c context.Context, file item.File, location string) (string, error) {
	userContext, ok := c.Value("user").(*user.ContextData)
	if !ok {
		log.Error(errors.New("context_retrieve_user_err"))
		return "", errors.New("context_retrieve_user_err")
	}
	mimeType := file.Header.Header.Get("Content-Type")
	ok = checkAllowedMimeTypes(mimeType)
	if !ok {
		return "", ErrBadParamInput
	}
	fileNames := strings.Split(file.Header.Filename, ".")
	uuid := uuid.NewV4()
	input := strings.NewReader(fileNames[0] + userContext.UserName + uuid.String())
	hash := sha256.New()
	if _, err := io.Copy(hash, input); err != nil {
		log.Fatal(err)
	}
	sum := hash.Sum(nil)
	file.Header.Filename = hex.EncodeToString(sum[:]) + generateFileName(mimeType)
	f, err := os.OpenFile("./assets/"+location+"/"+file.Header.Filename, os.O_WRONLY|os.O_CREATE, 0666)
	if err != nil {
		log.Error(err)
		return "", err
	}
	defer f.Close()
	_, err = io.Copy(f, file.Physical)
	if err != nil {
		log.Error(err)
		return "", err
	}
	fmt.Println("storing")
	return file.Header.Filename, err
}

// CheckAllowedMimeTypes checks if image is supported type
func checkAllowedMimeTypes(mimeType string) bool {
	for _, value := range allowedContentTypes {
		if mimeType == value {
			return true
		}
	}
	return false
}

// generateFileName generates file type extentions
func generateFileName(contentType string) string {

	var ext string

	switch contentType {
	case "image/jpeg":
		ext = ".jpg"
	case "image/jpg":
		ext = ".jpg"
	case "image/png":
		ext = ".png"
	case "image/gif":
		ext = ".gif"
	}

	return ext
}
