package main

import (
	"flag"
	"github.com/go-chi/chi/middleware"
	"log"
	"net/http"
	"os"
	"os/signal"

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
