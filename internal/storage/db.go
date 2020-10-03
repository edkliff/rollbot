package storage

import "github.com/edkliff/rollbot/internal/config"

type Storage interface {
	GetUser(int) (string, error)
	SetUser(int, string) error
	LoadUsers() error
}

const (
	sqlite  = "sqlite"
	file = "file"
)
func CreateStorage(conf config.Config) (Storage, error) {
	switch conf.DB.Kind {
	case sqlite:
		return ConnectSQLite(conf.DB)
	}
	return ConnectSQLite(conf.DB)
}