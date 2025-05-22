package repos

import (
	"context"
	"crypto/rand"
	"errors"
	"math/big"
	"quotes/models"
	"slices"
	"strconv"
)

type QuotesRepo struct {
	Quotes      []models.Quote
	lastQuoteID int
}

func NewQuotesRepo() *QuotesRepo {

	return &QuotesRepo{
		Quotes: []models.Quote{
			{ID: 1, Author: "Evgeny", Quote: "Люблю грозу в начале мая, и слушать Андрея Замая"},
			{ID: 2, Author: "Dmitry", Quote: "Цитата 2"},
		},
		lastQuoteID: 2,
	}

}

func (q *QuotesRepo) GetAll(ctx context.Context) ([]models.Quote, error) {

	if ctx.Err() != nil {
		return nil, ctx.Err()
	}

	return q.Quotes, nil
}

func (q *QuotesRepo) GetRandom(ctx context.Context) (*models.Quote, error) {

	if len(q.Quotes) == 1 {
		return &q.Quotes[0], nil
	}

	if len(q.Quotes) == 0 {
		return nil, errors.New("цитат нет")
	}

	if ctx.Err() != nil {
		return nil, ctx.Err()
	}

	lastNum := big.NewInt(int64(q.lastQuoteID))

	id, err := rand.Int(rand.Reader, lastNum)

	if err != nil {
		return nil, errors.New("не удалось получить случайное число" + err.Error())
	}

	return &q.Quotes[int(id.Int64())], nil

}

func (q *QuotesRepo) GetByAuthor(ctx context.Context, author string) ([]models.Quote, error) {

	if ctx.Err() != nil {
		return nil, ctx.Err()
	}

	var quotes []models.Quote

	for _, quote := range q.Quotes {
		if quote.Author == author {
			quotes = append(quotes, quote)
		}
	}

	if len(quotes) > 0 {
		return quotes, nil
	}

	return nil, errors.New("такого автора нет")
}

func (q *QuotesRepo) Create(ctx context.Context, quote *models.Quote) error {

	if ctx.Err() != nil {
		return ctx.Err()
	}

	q.lastQuoteID++

	q.Quotes = append(q.Quotes, models.Quote{
		ID:     q.lastQuoteID,
		Author: quote.Author,
		Quote:  quote.Quote},
	)

	return nil
}

func (q *QuotesRepo) Delete(ctx context.Context, id int) error {

	if ctx.Err() != nil {
		return ctx.Err()
	}

	if id > q.lastQuoteID {
		return errors.New("такого айди нет, последний айди в таблице: " + strconv.Itoa(q.lastQuoteID))
	}

	if id <= 0 {
		return errors.New("айди не может быть меньше или равен нулю")
	}

	id--

	q.Quotes = slices.Delete(q.Quotes, id, id+1)

	if q.lastQuoteID > 0 {
		q.lastQuoteID--
	}

	if len(q.Quotes) > 0 {
		for i := id; i < len(q.Quotes); i++ {
			q.Quotes[i].ID--
		}
	}

	return nil

}
