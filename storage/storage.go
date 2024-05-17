package storage

import (
	"FinancialBot/lib/e"
	"context"
	"crypto/sha1"
	"errors"
	"fmt"
	"io"
)

type Storage interface {
	Save(ctx context.Context, p *Page) error
	Remove(ctx context.Context, p *Page) error
	GetTotalIncome(ctx context.Context, username string) (float64, error)
	GetTotalExpense(ctx context.Context, username string) (float64, error)
}

var ErrNoSavedPages = errors.New("No saved pages")

type Page struct {
	UserName string
	Income   string
	Expenses string
	DateTime string
}

func (p Page) Hash() (string, error) {
	h := sha1.New()

	if _, err := io.WriteString(h, p.UserName); err != nil {
		return "", e.Wrap("can't calculate hash", err)
	}

	if _, err := io.WriteString(h, p.Income); err != nil {
		return "", e.Wrap("can't calculate hash", err)
	}

	if _, err := io.WriteString(h, p.Expenses); err != nil {
		return "", e.Wrap("can't calculate hash", err)
	}

	if _, err := io.WriteString(h, p.DateTime); err != nil {
		return "", e.Wrap("can't calculate hash", err)
	}

	return fmt.Sprintf("%x", h.Sum(nil)), nil
}
