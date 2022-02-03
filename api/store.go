package api

import (
	"database/sql"
	"fmt"

	_ "github.com/mattn/go-sqlite3"
)

type Store struct {
	Driver *sql.DB
	Path   string
}

func NewStore(path string) *Store {
	return &Store{
		Path: path,
	}
}

func (s *Store) Open() error {
	db, err := sql.Open("sqlite3", s.Path)
	if err != nil {
		return err
	}
	s.Driver = db
	return nil
}

func (s *Store) Close() error {
	return s.Driver.Close()
}

func (s *Store) Init(token string) error {
	if err := s.Open(); err != nil {
		return err
	}
	defer s.Close()

	if _, err := s.Driver.Exec(
		`CREATE TABLE IF NOT EXISTS stocks (
			id INTEGER PRIMARY KEY, 
			name VARCHAR UNIQUE, 
			ticker VARCHAR,
			price INTEGER, 
			cur VARCHAR
		)`, nil,
	); err != nil {
		return err
	}

	if _, err := s.Driver.Exec(
		"CREATE TABLE IF NOT EXISTS token (token VARCHAR UNIQUE)", nil,
	); err != nil {
		return err
	}

	if _, err := s.Driver.Exec(
		"INSERT INTO token (token) VALUES ($1)", token,
	); err != nil {
		return err
	}

	return nil
}

func (s *Store) GetToken() (string, error) {
	if err := s.Open(); err != nil {
		return "", err
	}
	defer s.Close()
	var token string
	row := s.Driver.QueryRow(
		"SELECT token FROM token LIMIT 1", nil,
	)
	if err := row.Scan(&token); err != nil {
		return "", err
	}
	return token, nil
}

func (s *Store) AddStock(stock Stock) error {
	if err := s.Open(); err != nil {
		return err
	}
	if _, err := s.Driver.Exec(
		"INSERT INTO stocks (name, ticker, price, cur) VALUES ($1, $2, $3, $4)",
		stock.Name, stock.Ticker, fmt.Sprintf("%.2f", stock.Price), stock.Currency,
	); err != nil {
		return err
	}
	defer s.Close()
	return nil
}

func (s *Store) UpdateStock(stock Stock) error {
	if err := s.Open(); err != nil {
		return err
	}
	if _, err := s.Driver.Exec(
		"UPDATE stocks SET price = $1 WHERE name = $2",
		fmt.Sprintf("%.2f", stock.Price), stock.Name,
	); err != nil {
		return err
	}
	defer s.Close()
	return nil
}

func (s *Store) GetStockByTicker(ticker string) (Stock, error) {
	if err := s.Open(); err != nil {
		return Stock{}, err
	}
	var stock Stock
	row := s.Driver.QueryRow(
		"SELECT name, ticker, price, cur FROM stocks WHERE ticker = $1", ticker,
	)
	if err := row.Scan(&stock.Name, &stock.Ticker, &stock.Price, &stock.Currency); err != nil {
		return Stock{}, err
	}
	return stock, nil
}

func (s *Store) GetStockRand() (Stock, error) {
	if err := s.Open(); err != nil {
		return Stock{}, err
	}
	var stock Stock
	row := s.Driver.QueryRow(
		"SELECT name, ticker, price, cur FROM stocks ORDER BY RANDOM() LIMIT 1;",
	)
	if err := row.Scan(&stock.Name, &stock.Ticker, &stock.Price, &stock.Currency); err != nil {
		return Stock{}, err
	}
	return stock, nil
}

func (s *Store) DelStockByTicker(ticker string) error {
	if err := s.Open(); err != nil {
		return err
	}
	if _, err := s.Driver.Exec(
		"DELETE FROM stocks WHERE ticker = $1", ticker,
	); err != nil {
		return err
	}

	return nil
}
