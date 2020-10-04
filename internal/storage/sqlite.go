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
	err = sqlConn.CreateDB()
	if err != nil {
		return nil, err
	}
	err = sqlConn.LoadUsers()
	if err != nil {
		return nil, err
	}
	return &sqlConn, nil
}

func (s *SQLiteConnection)  CreateDB() error  {
	_, err := s.Database.Exec(`CREATE TABLE IF NOT EXISTS users (
    id INTEGER PRIMARY KEY UNIQUE,
    username TEXT
)`)
	if err != nil {
		return err
	}
	_, err = s.Database.Exec(`CREATE TABLE IF NOT EXISTS logs (
    id INTEGER PRIMARY KEY AUTOINCREMENT UNIQUE,
    username int,
    command TEXT,
    result TEXT,
    date integer
)`)
	if err != nil {
		return err
	}
	return nil
}

func (s *SQLiteConnection) GetUser(userId int) (string, error) {
	user, ok := s.Users.GetUser(userId)
	if !ok {
		return "", errors.New("unknown user")
	}
	return user, nil
}

func (s *SQLiteConnection) SetUser(userID int, username string) error {
	 s.Users.SetUser(userID, username)
	 _, err:= s.Database.Exec(`INSERT INTO users (id, username) VALUES ($1, $2)`, userID, username)
	 if err != nil {
	 	return err
	 }
	return nil
}

func (s *SQLiteConnection) LoadUsers() error {
	rows, err := s.Database.Query(`SELECT id, username FROM users`)
	if err != nil {
		return err
	}
	defer rows.Close()
	for rows.Next() {
		id := 0
		name := ""
		err := rows.Scan(&id, &name)
		if err != nil {
			return err
		}
		s.Users.SetUser(id, name)
	}
	return nil
}
