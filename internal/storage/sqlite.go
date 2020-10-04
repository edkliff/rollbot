package storage

import (
	"database/sql"
	"errors"
	"fmt"
	"github.com/edkliff/rollbot/internal/config"
	_ "github.com/mattn/go-sqlite3"
	"time"
)

type SQLiteConnection struct {
	Database *sql.DB
	Users    UserCache
}

func ConnectSQLite(conf config.DBConfig) (*SQLiteConnection, error) {
	dbfile := "rollbot.db"
	if conf.Filename != "" {
		dbfile = conf.Filename
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

func (s *SQLiteConnection) CreateDB() error {
	_, err := s.Database.Exec(`CREATE TABLE IF NOT EXISTS users (
    id INTEGER PRIMARY KEY UNIQUE,
    username TEXT
)`)
	if err != nil {
		return err
	}
	_, err = s.Database.Exec(`CREATE TABLE IF NOT EXISTS logs (
    id INTEGER PRIMARY KEY AUTOINCREMENT UNIQUE,
    user_id int,
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
	_, err := s.Database.Exec(`INSERT INTO users (id, username) VALUES ($1, $2)`, userID, username)
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

func (s *SQLiteConnection) UsersList() string {
	return fmt.Sprintf("%v", s.Users.users)
}

type User struct {
	ID int
	Username string
	Count int
}

type UsersList struct {
	Users []User
}

func (s *SQLiteConnection) GetUsers() (*UsersList, error)  {
	users := make([]User,0)
	q := `SELECT
		u.id, u.username, count(l.id)
		FROM users u 
		LEFT OUTER JOIN logs l on l.user_id  = u.id
		GROUP BY 	u.id, u.username
		ORDER BY count(l.id) DESC `
	rows, err := s.Database.Query(q)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	for rows.Next() {
		u := User{}
		err := rows.Scan(&u.ID, &u.Username, &u.Count)
		if err != nil {
			return nil, err
		}
		users = append(users, u)
	}
	l := UsersList{Users:users}
	return &l, nil
}

func (s *SQLiteConnection) WriteTask(original string, response string, user int) error {
	date := time.Now().Unix()
	q := `INSERT INTO logs (user_id, command, result, date) VALUES ($1, $2, $3, $4)`
	_, err := s.Database.Exec(q, user, original, response, date)
	if err != nil {
		return err
	}
	return nil
}

type Result struct {
	UserID int
	Username string
	Command string
	Result string
	Date string
}

type ResultsList struct {
	Results []Result
}


func (s *SQLiteConnection) GetLogs(userid int) (*ResultsList, error)  {
	results := make([]Result,0)
	q := `SELECT u.username, u.id, l.command, l.result, l.date FROM logs l 
			JOIN users u on l.user_id  = u.id `
	args := make([]interface{}, 0)
	if userid != 0 {
		q += ` WHERE u.id = $1 `
		args = append(args, userid)
	}
	q += `ORDER BY l.date DESC LIMIT 100`
	rows, err := s.Database.Query(q, args...)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	for rows.Next() {
		r := Result{}
		d := int64(0)
		err := rows.Scan(&r.Username, &r.UserID, &r.Command, &r.Result, &d)
		if err != nil {
			return nil, err
		}
		date := time.Unix(d, 0).String()
		r.Date = date
		results = append(results, r)
	}
	rl := ResultsList{Results:results}
	return &rl, nil
}
