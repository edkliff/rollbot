package app

import (
	"github.com/edkliff/rollbot/internal/config"
	"github.com/edkliff/rollbot/internal/generator"
	"github.com/edkliff/rollbot/internal/storage"
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

