package models

import "context"

type QuoteReader interface {
	GetAll(ctx context.Context) ([]Quote, error)
	GetRandom(ctx context.Context) (*Quote, error)
	GetByAuthor(ctx context.Context, author string) ([]Quote, error)
}

type QuoteWriter interface {
	Create(ctx context.Context, quote *Quote) error
}

type QuoteDeleter interface {
	Delete(ctx context.Context, id int) error
}

type QuoteRepo interface {
	QuoteReader
	QuoteWriter
	QuoteDeleter
}

type Quote struct {
	ID     int
	Author string
	Quote  string
}
