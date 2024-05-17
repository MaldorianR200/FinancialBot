package sqlite

import (
	"FinancialBot/storage"
	"context"
	"database/sql"
	"fmt"
	_ "github.com/mattn/go-sqlite3"
)

type Storage struct {
	db *sql.DB
}

// New creates new SQLite storage.
func New(path string) (*Storage, error) {
	db, err := sql.Open("sqlite3", path)
	if err != nil {
		return nil, fmt.Errorf("can't open database: %w", err)
	}

	if err := db.Ping(); err != nil {
		return nil, fmt.Errorf("can't connect database: %w", err)
	}
	return &Storage{db: db}, nil
}

// Save сохраняет данные в хранилище.
func (s *Storage) Save(ctx context.Context, page *storage.Page) error {
	q := `INSERT INTO pages (user_name, income, expenses, datetime) VALUES (?, ?, ?, ?)`

	if _, err := s.db.ExecContext(ctx, q, page.UserName, page.Income, page.Expenses, page.DateTime); err != nil {
		return fmt.Errorf("can't save page: %w", err)
	}

	return nil
}

// GetTotalIncome получает общий доход пользователя.
func (s *Storage) GetTotalIncome(ctx context.Context, username string) (float64, error) {
	q := `SELECT SUM(income) FROM pages WHERE user_name = ?`
	var totalIncome float64

	row := s.db.QueryRowContext(ctx, q, username)
	if err := row.Scan(&totalIncome); err != nil {
		return 0, fmt.Errorf("can't get total income: %w", err)
	}

	return totalIncome, nil
}

// GetTotalExpense - получает общие расходы пользователя.
func (s *Storage) GetTotalExpense(ctx context.Context, username string) (float64, error) {
	q := `SELECT SUM(expenses) FROM pages WHERE user_name = ?`
	var totalExpense float64

	row := s.db.QueryRowContext(ctx, q, username)
	if err := row.Scan(&totalExpense); err != nil {
		return 0, fmt.Errorf("can't get total expense: %w", err)
	}

	return totalExpense, nil
}

// Remove -
func (s *Storage) Remove(ctx context.Context, page *storage.Page) error {
	q := `DELETE FROM pages WHERE user_name = ? AND income = ? AND expenses = ? AND datetime = ? `
	//q := `DELETE FROM pages WHERE user_name = ? AND datetime = ?`
	if _, err := s.db.ExecContext(ctx, q, page.UserName, page.DateTime); err != nil {
		return fmt.Errorf("can't remove page: %w", err)
	}

	return nil
}

func (s *Storage) Init(ctx context.Context) error {
	q := `CREATE TABLE IF NOT EXISTS pages (
    id INTEGER PRIMARY KEY AUTOINCREMENT,
    user_name TEXT, 
    income INT, 
    expenses int, 
    datetime VARCHAR
    )`

	if _, err := s.db.ExecContext(ctx, q); err != nil {
		return fmt.Errorf("can't create table: %w", err)
	}

	return nil
}
