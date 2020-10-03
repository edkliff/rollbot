package storage
import (
	"errors"
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

func (s *SQLiteConnection) GetUser(userId int) (string, error) {
	user, ok := s.Users.GetUser(userId)
	if !ok {
		return "", errors.New("unknown user")
	}
	return user, nil
}

func (s *SQLiteConnection) SetUser( userID int, username string) error {
	 s.Users.SetUser(userID, username)
	return nil
}

func (s *SQLiteConnection) LoadUsers() error {
	return nil
}
