package app

import (
	"github.com/edkliff/rollbot/internal/config"
	"github.com/edkliff/rollbot/internal/storage"
	"github.com/edkliff/rollbot/internal/generator"
)

type RollBot struct {
	Config    config.Config
	DB        storage.Storage
	Generator *generator.Generator
}

func CreateRollBot(conf config.Config, store storage.Storage) *RollBot {
	return &RollBot{
		Config:    conf,
		DB:        store,
		Generator: generator.InitGenerator(),
	}
}
