package routes

import (
	_itemRepo "go_boilerplate/item/repository"
	_itemUseCase "go_boilerplate/item/usecase"
	_lib "go_boilerplate/lib"
	"go_boilerplate/middleware"
	_userRepo "go_boilerplate/user/repository"
	_userUseCase "go_boilerplate/user/usecase"
	"os"
	"strconv"
	"time"

	"github.com/Masterminds/squirrel"
	"github.com/go-chi/chi"
	"github.com/go-chi/cors"
	"github.com/go-redis/redis"
	"github.com/jmoiron/sqlx"
)

var log = _lib.GetLogger()

func InitRoutes(db *sqlx.DB, sb squirrel.StatementBuilderType, rs *redis.Client) *chi.Mux {
	log.Infoln("Setting up routes...")
	// router init
	tc, _ := strconv.Atoi(os.Getenv("APP_CTX_TIMEOUT"))
	timeoutContext := time.Duration(tc) * time.Second
	r := chi.NewRouter()
	// Basic CORS
	// for more ideas, see: https://developer.github.com/v3/#cross-origin-resource-sharing
	cors := cors.New(cors.Options{
		// AllowedOrigins: []string{"https://foo.com"}, // Use this to allow specific origin hosts
		AllowedOrigins: []string{"*"},
		// AllowOriginFunc:  func(r *http.Request, origin string) bool { return true },
		AllowedMethods:   []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
		AllowedHeaders:   []string{"Accept", "Authorization", "Content-Type", "X-CSRF-Token"},
		ExposedHeaders:   []string{"Link"},
		AllowCredentials: true,
		MaxAge:           300, // Maximum value not ignored by any of major browsers
	})

	// apply middlewares
	log.Infoln("Applying HTTP middlewares...")
	r.Use(middleware.RequestLogger)
	r.Use(cors.Handler)

	// inject domains with their dependencies and setup their routes
	log.Infoln("Injecting services dependencies and setting up their routes...")
	// user routes
	userRepo := _userRepo.NewUserRepository(sb, db, rs)
	userUse := _userUseCase.NewUserUseCase(userRepo, timeoutContext)
	UserHttpRouter(r, userUse)

	// item routes
	itemRepo := _itemRepo.NewItemRepository(sb, db)
	ju := _itemUseCase.NewItemUseCase(itemRepo, userRepo, timeoutContext)
	ItemHttpRouter(r, ju)
	return r
}
