package repos

import (
	"context"
	"quotes/models"
	"testing"
)

func TestNewQuotesRepo(t *testing.T) {
	repo := NewQuotesRepo()

	if len(repo.Quotes) != 2 {
		t.Errorf("Ожидалось 2 начальные цитаты, получено %d", len(repo.Quotes))
	}

	if repo.lastQuoteID != 2 {
		t.Errorf("Ожидалось lastQuoteID равным 2, получено %d", repo.lastQuoteID)
	}
}

func TestGetAll(t *testing.T) {
	tests := []struct {
		name    string
		ctx     context.Context
		wantLen int
		wantErr bool
	}{
		{
			name:    "успешный вызов",
			ctx:     context.Background(),
			wantLen: 2,
			wantErr: false,
		},
		{
			name:    "отмененный контекст",
			ctx:     cancelledContext(),
			wantLen: 0,
			wantErr: true,
		},
	}

	repo := NewQuotesRepo()

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			quotes, err := repo.GetAll(tt.ctx)

			if (err != nil) != tt.wantErr {
				t.Errorf("GetAll() ошибка = %v, ожидалась ошибка %v", err, tt.wantErr)
				return
			}

			if !tt.wantErr && len(quotes) != tt.wantLen {
				t.Errorf("GetAll() = %v, ожидалось %v", len(quotes), tt.wantLen)
			}
		})
	}
}

func TestGetRandom(t *testing.T) {
	tests := []struct {
		name    string
		ctx     context.Context
		repo    *QuotesRepo
		wantErr bool
	}{
		{
			name:    "успешный вызов",
			ctx:     context.Background(),
			repo:    NewQuotesRepo(),
			wantErr: false,
		},
		{
			name:    "отмененный контекст",
			ctx:     cancelledContext(),
			repo:    NewQuotesRepo(),
			wantErr: true,
		},
		{
			name: "нет цитат",
			ctx:  context.Background(),
			repo: &QuotesRepo{
				Quotes:      []models.Quote{},
				lastQuoteID: 0,
			},
			wantErr: true,
		},
		{
			name: "одна цитата",
			ctx:  context.Background(),
			repo: &QuotesRepo{
				Quotes:      []models.Quote{{ID: 1, Author: "Test", Quote: "Test"}},
				lastQuoteID: 1,
			},
			wantErr: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			_, err := tt.repo.GetRandom(tt.ctx)

			if (err != nil) != tt.wantErr {
				t.Errorf("GetRandom() ошибка = %v, ожидалась ошибка %v", err, tt.wantErr)
			}
		})
	}
}

func TestGetByAuthor(t *testing.T) {
	repo := NewQuotesRepo()

	tests := []struct {
		name    string
		ctx     context.Context
		author  string
		wantLen int
		wantErr bool
		errMsg  string
	}{
		{
			name:    "успешный поиск",
			ctx:     context.Background(),
			author:  "Evgeny",
			wantLen: 1,
			wantErr: false,
		},
		{
			name:    "автор не найден",
			ctx:     context.Background(),
			author:  "Unknown",
			wantLen: 0,
			wantErr: true,
			errMsg:  "такого автора нет",
		},
		{
			name:    "отмененный контекст",
			ctx:     cancelledContext(),
			author:  "Evgeny",
			wantLen: 0,
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			quotes, err := repo.GetByAuthor(tt.ctx, tt.author)

			if (err != nil) != tt.wantErr {
				t.Errorf("GetByAuthor() ошибка = %v, ожидалась ошибка %v", err, tt.wantErr)
				return
			}

			if tt.wantErr && tt.errMsg != "" && err.Error() != tt.errMsg {
				t.Errorf("GetByAuthor() сообщение ошибки = '%v', ожидалось '%v'", err.Error(), tt.errMsg)
			}

			if !tt.wantErr && len(quotes) != tt.wantLen {
				t.Errorf("GetByAuthor() количество цитат = %v, ожидалось %v", len(quotes), tt.wantLen)
			}
		})
	}
}

func TestCreate(t *testing.T) {
	repo := NewQuotesRepo()
	initialCount := len(repo.Quotes)

	tests := []struct {
		name    string
		ctx     context.Context
		quote   *models.Quote
		wantErr bool
	}{
		{
			name:    "успешное создание",
			ctx:     context.Background(),
			quote:   &models.Quote{Author: "New", Quote: "New quote"},
			wantErr: false,
		},
		{
			name:    "отмененный контекст",
			ctx:     cancelledContext(),
			quote:   &models.Quote{Author: "New", Quote: "New quote"},
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := repo.Create(tt.ctx, tt.quote)

			if (err != nil) != tt.wantErr {
				t.Errorf("Create() ошибка = %v, ожидалась ошибка %v", err, tt.wantErr)
				return
			}

			if !tt.wantErr {
				if len(repo.Quotes) != initialCount+1 {
					t.Errorf("После Create() количество цитат = %v, ожидалось %v", len(repo.Quotes), initialCount+1)
				}
				initialCount = len(repo.Quotes)
			}
		})
	}
}

func TestDelete(t *testing.T) {
	tests := []struct {
		name    string
		ctx     context.Context
		id      int
		repo    *QuotesRepo
		wantErr bool
		errMsg  string
	}{
		{
			name:    "успешное удаление",
			ctx:     context.Background(),
			id:      1,
			repo:    NewQuotesRepo(),
			wantErr: false,
		},
		{
			name:    "отмененный контекст",
			ctx:     cancelledContext(),
			id:      1,
			repo:    NewQuotesRepo(),
			wantErr: true,
		},
		{
			name:    "несуществующий ID",
			ctx:     context.Background(),
			id:      100,
			repo:    NewQuotesRepo(),
			wantErr: true,
			errMsg:  "такого айди нет, последний айди в таблице: 2",
		},
		{
			name:    "ID <= 0",
			ctx:     context.Background(),
			id:      0,
			repo:    NewQuotesRepo(),
			wantErr: true,
			errMsg:  "айди не может быть меньше или равен нулю",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			initialCount := len(tt.repo.Quotes)
			err := tt.repo.Delete(tt.ctx, tt.id)

			if (err != nil) != tt.wantErr {
				t.Errorf("Delete() ошибка = %v, ожидалась ошибка %v", err, tt.wantErr)
				return
			}

			if tt.wantErr && tt.errMsg != "" && err.Error() != tt.errMsg {
				t.Errorf("Delete() сообщение ошибки = '%v', ожидалось '%v'", err.Error(), tt.errMsg)
			}

			if !tt.wantErr {
				if len(tt.repo.Quotes) != initialCount-1 {
					t.Errorf("После Delete() количество цитат = %v, ожидалось %v", len(tt.repo.Quotes), initialCount-1)
				}
			}
		})
	}
}

func cancelledContext() context.Context {
	ctx, cancel := context.WithCancel(context.Background())
	cancel()
	return ctx
}
