package mysql

import (
	"database/sql"
	"fmt"
	_ "github.com/go-sql-driver/mysql"
	"os"
)

type Storage struct {
	db *sql.DB
}

func New(addressDB string, login string, pass string, nameDB string) (*Storage, error) {
	os.Setenv("CGO_ENABLED", "1")
	const op = "storage.mysql.New"

	db, err := sql.Open("mysql", fmt.Sprintf("%s:%s@tcp(%s)/%s", login, pass, addressDB, nameDB))
	if err != nil {
		return nil, fmt.Errorf("%s: %w", op, err)
	}

	db.Query("SET GLOBAL sql_mode=''")

	// Создаем таблицу, если ее еще нет
	insert, err := db.Query(`
	CREATE TABLE IF NOT EXISTS url(
	   id INTEGER PRIMARY KEY AUTO_INCREMENT,
	   alias VARCHAR(200) NOT NULL UNIQUE,
	   url TEXT NOT NULL);
	`)
	if err != nil {
		return nil, fmt.Errorf("%s: %w", op, err)
	}
	_ = insert

	return &Storage{db: db}, nil
}

func (s *Storage) SaveURL(urlToSave string, alias string) (int64, error) {
	const op = "storage.mysql.SaveURL"
	stmt, err := s.db.Prepare("INSERT INTO url(url, alias) VALUES(?, ?);")
	if err != nil {
		return 0, fmt.Errorf("%s, %w", op, err)
	}
	res, err := stmt.Exec(urlToSave, alias)
	if err != nil {
		return 0, fmt.Errorf("%s, %w", op, err)
	}
	id, err := res.LastInsertId()
	fmt.Println("id: ", id)
	return id, nil
}

func (s *Storage) GetUrl(alias string) (string, error) {
	const op = "storage.mysql.GetUrl"
	stmt, err := s.db.Prepare("SELECT url from url WHERE alias = ?")
	if err != nil {
		return "", fmt.Errorf("%s, %w", op, err)
	}

	var resUrl string
	err = stmt.QueryRow(alias).Scan(&resUrl)
	if err != nil {
		return "", fmt.Errorf("%s, %w", op, err)
	}

	return resUrl, nil
}
