package main

import (
	"fmt"
	"go_boilerplate/config"
	_itemHttpDelivery "go_boilerplate/item/delivery/http"
	_itemRepo "go_boilerplate/item/repository"
	_itemUseCase "go_boilerplate/item/usecase"
	_lib "go_boilerplate/lib"
	_userHttpDelivery "go_boilerplate/user/delivery/http"
	_userRepo "go_boilerplate/user/repository"
	_userUseCase "go_boilerplate/user/usecase"
	"net/http"
	"os"
	"path/filepath"
	"strings"
	"time"

	"github.com/go-chi/chi"
	"github.com/go-chi/cors"
	"github.com/go-redis/redis"
	_ "github.com/go-sql-driver/mysql"
	"github.com/jmoiron/sqlx"
	"github.com/sirupsen/logrus"
)

func main() {
	log := _lib.GetLogger()

	// db startup
	v, err := config.ReadConfig("db")
	if err != nil {
		log.Fatalf("Error when reading config: %v", err)
	}
	db, err := sqlx.Connect("mysql", fmt.Sprintf("%s:%s@tcp(%s:%s)/%s?parseTime=%s",
		v.Get("Username"), v.Get("Password"), v.Get("Host"), v.Get("Port"), v.Get("Database"), v.Get("ParseTime")))
	if err != nil {
		log.Fatalf("Failed connecting to the database: %v", err)
	}
	defer db.Close()

	// redis
	redisConn := redis.NewClient(&redis.Options{
		Addr:     "localhost:6379",
		Password: "", // no password set
		DB:       0,  // use default DB
	})
	_, err = redisConn.Ping().Result()
	if err != nil {
		log.Fatalf("Redis connection failed: %v\n", err)
	}
	defer redisConn.Close()

	// business logic init
	timeoutContext := 3000 * time.Second
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
	r.Use(cors.Handler)

	// inject domains
	userRepo := _userRepo.NewUserRepository(db, redisConn)
	userUse := _userUseCase.NewUserUseCase(userRepo, timeoutContext)
	_userHttpDelivery.UserHttpRouter(r, userUse)
	itemRepo := _itemRepo.NewItemRepository(db)
	ju := _itemUseCase.NewItemUseCase(itemRepo, userRepo, timeoutContext)
	_itemHttpDelivery.ItemHttpRouter(r, ju)

	//static files
	workDir, _ := os.Getwd()
	filesDir := filepath.Join(workDir, "assets")
	FileServer(r, "/public", http.Dir(filesDir), log)
	log.Info("App started at port 4000")
	// start server
	http.ListenAndServe(":4000", r)
}

// FileServer conveniently sets up a http.FileServer handler to serve
// static files from a http.FileSystem.
func FileServer(r chi.Router, path string, root http.FileSystem, log *logrus.Logger) {
	if strings.ContainsAny(path, "{}*") {
		log.Fatalf("FileServer does not permit URL parameters, path: %v", path)
	}

	fs := http.StripPrefix(path, http.FileServer(root))

	if path != "/" && path[len(path)-1] != '/' {
		r.Get(path, http.RedirectHandler(path+"/", 301).ServeHTTP)
		path += "/"
	}
	path += "*"

	r.Get(path, http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		fs.ServeHTTP(w, r)
	}))
}
