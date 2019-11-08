package http

import (
	"context"
	"go_boilerplate/item"
	"go_boilerplate/lib"
	"go_boilerplate/middleware"
	"go_boilerplate/models"
	"net/http"
	"strconv"

	_lib "go_boilerplate/lib"

	"github.com/asaskevich/govalidator"
	"github.com/go-chi/chi"
)

var log = _lib.GetLogger()

type NewHttpItemHandler struct {
	ItemUseCase item.UseCase
}

func ItemHttpRouter(router *chi.Mux, UseCase item.UseCase) {
	handler := &NewHttpItemHandler{
		ItemUseCase: UseCase,
	}
	r := chi.NewRouter()
	r.Get("/{id}", handler.GetByID)
	r.Group(func(r chi.Router) {
		r.Use(middleware.AccessTokenVerifier)
		r.Delete("/{id}", handler.Delete)
		r.Post("/", handler.Store)
	})
	r.Get("/", handler.FetchItem)

	router.Mount("/items", r)
}

func (i *NewHttpItemHandler) FetchItem(w http.ResponseWriter, r *http.Request) {
	offset := chi.URLParam(r, "offset")
	num, _ := strconv.Atoi(offset)
	ctx := r.Context()
	if ctx == nil {
		ctx = context.Background()
	}
	listItems, err := i.ItemUseCase.Fetch(ctx, int64(num))
	if err != nil {
		lib.RespondJSON(w, lib.GetStatusCode(err), nil, err.Error())
		return
	}
	lib.RespondJSON(w, http.StatusOK, listItems, "")
}

func (i *NewHttpItemHandler) GetByID(w http.ResponseWriter, r *http.Request) {

	idP, err := strconv.Atoi(chi.URLParam(r, "id"))
	id := int64(idP)

	ctx := r.Context()
	if ctx == nil {
		ctx = context.Background()
	}

	itemRow, err := i.ItemUseCase.GetByID(ctx, id)

	if err != nil {
		lib.RespondJSON(w, lib.GetStatusCode(err), nil, err.Error())
		return
	}
	lib.RespondJSON(w, http.StatusOK, itemRow, "")
}

func (i *NewHttpItemHandler) Store(w http.ResponseWriter, r *http.Request) {
	var itemRow models.Item

	ctx := r.Context()
	if ctx == nil {
		ctx = context.Background()
	}

	err := r.ParseMultipartForm(32 << 20)
	if err != nil {
		log.Error(err.Error())
		lib.RespondJSON(w, http.StatusUnprocessableEntity, nil, _lib.ErrBadParamInput.Error())
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
		f, _ := fh.Open()
		images = append(images, item.File{
			Physical: f,
			Header:   fh,
		})
	}
	if err != nil {
		log.Error(err.Error())
		lib.RespondJSON(w, http.StatusUnprocessableEntity, nil, _lib.ErrBadParamInput.Error())
		return
	}
	if ok, err := govalidator.ValidateStruct(&itemRow); !ok {
		lib.RespondJSON(w, http.StatusBadRequest, nil, err.Error())
		return
	}

	id, err := i.ItemUseCase.Store(ctx, &itemRow, images)

	if err != nil {
		lib.RespondJSON(w, lib.GetStatusCode(err), nil, err.Error())
		return
	}
	itemRow.ID = id
	lib.RespondJSON(w, http.StatusCreated, itemRow, "")
}

func (i *NewHttpItemHandler) Delete(w http.ResponseWriter, r *http.Request) {
	idP, err := strconv.Atoi(chi.URLParam(r, "id"))
	id := int64(idP)
	ctx := r.Context()
	if ctx == nil {
		ctx = context.Background()
	}

	err = i.ItemUseCase.Delete(ctx, id)

	if err != nil {
		lib.RespondJSON(w, lib.GetStatusCode(err), nil, err.Error())
		return
	}
	lib.RespondJSON(w, http.StatusNoContent, nil, "")
}
