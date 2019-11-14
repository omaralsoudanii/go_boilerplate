package http

import (
	"context"
	"go_boilerplate/item"
	lib "go_boilerplate/lib"
	"go_boilerplate/models"
	"net/http"
	"strconv"

	"github.com/asaskevich/govalidator"
	"github.com/go-chi/chi"
)

var log = lib.GetLogger()

type NewHttpItemHandler struct {
	ItemUseCase item.UseCase
}

func (i *NewHttpItemHandler) GetAllItem(w http.ResponseWriter, r *http.Request) {
	offset := chi.URLParam(r, "offset")
	num, _ := strconv.Atoi(offset)
	ctx := r.Context()
	if ctx == nil {
		ctx = context.Background()
	}
	listItems, err := i.ItemUseCase.GetAll(ctx, uint(num))
	if err != nil {
		log.Errorln(err)
		lib.RespondJSON(w, lib.GetStatusCode(err), nil, err)
		return
	}
	lib.RespondJSON(w, http.StatusOK, listItems, nil)
}

func (i *NewHttpItemHandler) GetByID(w http.ResponseWriter, r *http.Request) {

	idP, err := strconv.Atoi(chi.URLParam(r, "id"))
	id := uint(idP)

	ctx := r.Context()
	if ctx == nil {
		ctx = context.Background()
	}

	itemRow, err := i.ItemUseCase.GetByID(ctx, id)

	if err != nil {
		log.Errorln(err)
		lib.RespondJSON(w, lib.GetStatusCode(err), nil, err)
		return
	}
	lib.RespondJSON(w, http.StatusOK, itemRow, nil)
}

func (i *NewHttpItemHandler) Store(w http.ResponseWriter, r *http.Request) {
	var itemRow models.Item

	ctx := r.Context()
	if ctx == nil {
		ctx = context.Background()
	}

	err := r.ParseMultipartForm(32 << 20)
	if err != nil {
		log.Errorln(err)
		lib.RespondJSON(w, http.StatusUnprocessableEntity, nil, lib.ErrBadParamInput)
		return
	}
	itemRow.Title = r.FormValue("title")
	itemRow.Description = r.FormValue("description")
	price, err := strconv.ParseFloat(r.FormValue("price"), 64)
	CategoryID, err := strconv.Atoi(r.FormValue("category_id"))
	itemRow.Price = price
	itemRow.CategoryID = CategoryID

	files := r.MultipartForm.File["images"]
	images := []item.File{}
	for _, fh := range files {
		f, err := fh.Open()
		if err != nil {
			log.Errorln(err)
			lib.RespondJSON(w, http.StatusUnprocessableEntity, nil, lib.ErrBadParamInput)
			return
		}
		images = append(images, item.File{
			Physical: f,
			Header:   fh,
		})
	}
	if ok, err := govalidator.ValidateStruct(&itemRow); !ok {
		log.Warning(err)
		lib.RespondJSON(w, http.StatusBadRequest, nil, err)
		return
	}

	id, err := i.ItemUseCase.Store(ctx, &itemRow, images)

	if err != nil {
		log.Errorln(err)
		lib.RespondJSON(w, lib.GetStatusCode(err), nil, err)
		return
	}
	itemRow.ID = id
	lib.RespondJSON(w, http.StatusCreated, itemRow, nil)
}

func (i *NewHttpItemHandler) Delete(w http.ResponseWriter, r *http.Request) {
	idP, err := strconv.Atoi(chi.URLParam(r, "id"))
	id := uint(idP)
	ctx := r.Context()
	if ctx == nil {
		ctx = context.Background()
	}

	err = i.ItemUseCase.Delete(ctx, id)

	if err != nil {
		lib.RespondJSON(w, lib.GetStatusCode(err), nil, err)
		return
	}
	lib.RespondJSON(w, http.StatusNoContent, nil, nil)
}
