package storage
import (
	"github.com/edkliff/rollbot/internal/config"
	"database/sql"
	_ "github.com/mattn/go-sqlite3"
)
type SQLiteConnection struct {
	Database *sql.DB
	Users UserCache
}

func ConnectSQLite(conf config.DBConfig) (*SQLiteConnection, error)  {
	dbfile := "rollbot.db"
	if conf.Filename != "" {
		dbfile =  conf.Filename
	}
	db, err := sql.Open("sqlite3", dbfile)
	if err != nil {
		return nil, err
	}
	sqlConn := SQLiteConnection{
		Database: db,
		Users:    NewUserCache(),
	}
	return &sqlConn, nil
}

func (s *SQLiteConnection)  CreateDB() error  {
	return nil
}

func (s *SQLiteConnection) GetUser(username string) (string, bool) {
	return s.Users.GetUser(username)
}

func (s *SQLiteConnection) SetUser(username string, userID string) error {
	 s.Users.SetUser(username, userID)
	return nil
}

func (s *SQLiteConnection) LoadUsers() error {
	return nil
}
