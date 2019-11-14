package main

import (
	"go_boilerplate/database"
	_lib "go_boilerplate/lib"
	"go_boilerplate/redis"
	"go_boilerplate/routes"
	"net/http"
	"os"
	"path/filepath"
	"strings"

	"github.com/go-chi/chi"
	_ "github.com/go-sql-driver/mysql"

	"github.com/sirupsen/logrus"
)

func main() {
	// log
	log := _lib.GetLogger()

	// set envs
	_lib.GetENV()

	// db startup
	db, sb := database.GetInstance(log)
	defer db.Close()

	// redis
	rs := redis.GetInstance(log)
	defer rs.Close()
	// init router
	r := routes.InitRoutes(db, sb, rs)

	//static files
	workDir, _ := os.Getwd()
	filesDir := filepath.Join(workDir, "assets")
	FileServer(r, "/public", http.Dir(filesDir), log)

	// start server
	log.Infoln("App started at port 4000")
	http.ListenAndServe(":4000", r)
}

// FileServer conveniently sets up a http.FileServer handler to serve
// static files from a http.FileSystem.
func FileServer(r *chi.Mux, path string, root http.FileSystem, log *logrus.Logger) {
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
