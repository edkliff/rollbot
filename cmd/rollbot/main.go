package main

import (
	"flag"
	"github.com/go-chi/chi/middleware"
	"log"
	"net/http"
	"os"
	"os/signal"
	"strings"

	"github.com/edkliff/rollbot/internal/app"
	"github.com/edkliff/rollbot/internal/config"
	"github.com/edkliff/rollbot/internal/storage"

	"github.com/go-chi/chi"
)

func main() {
	configFileName := flag.String("c", "c.yaml", "specify path to a c-example.yaml")
	flag.Parse()

	conf, err := config.LoadConfig(*configFileName)
	if err != nil {
		log.Fatal(err)
	}
	log.Printf("Config: %+v\n", conf)
	store, err := storage.CreateStorage(*conf)
	if err != nil {
		log.Fatal(err)
	}
	rollbot := app.CreateRollBot(*conf, store)

	mux := chi.NewRouter()

	imgs := "./img"
	fs := http.Dir(imgs)
	FileServer(mux, "/img", fs)

	mux.With(middleware.SetHeader("Content-Type", "text/html")).Get("/home", rollbot.Homepage)
	mux.With(middleware.SetHeader("Content-Type", "text/json")).
		Post("/vk", rollbot.VKHandle)
	// mux.Get("/users", rollbot.GetUsers)
	// mux.Get("/history", rollbot.GetHistory)
	// mux.Get("/history/{userId}", rollbot.GetUserHistory)


	server := http.Server{
		Handler: mux,
		Addr:    conf.Server,
	}

	err = server.ListenAndServe()
	if err != nil {
		log.Fatal(err)
	}
	sgnl := make(chan os.Signal)
	signal.Notify(sgnl, os.Interrupt, os.Kill)
	<-sgnl
}

func FileServer(r chi.Router, path string, root http.FileSystem) {
	if strings.ContainsAny(path, "{}*") {
		panic("FileServer does not permit any URL parameters.")
	}

	if path != "/" && path[len(path)-1] != '/' {
		r.Get(path, http.RedirectHandler(path+"/", 301).ServeHTTP)
		path += "/"
	}
	path += "*"

	r.Get(path, func(w http.ResponseWriter, r *http.Request) {
		rctx := chi.RouteContext(r.Context())
		pathPrefix := strings.TrimSuffix(rctx.RoutePattern(), "/*")
		fs := http.StripPrefix(pathPrefix, http.FileServer(root))
		fs.ServeHTTP(w, r)
	})
}