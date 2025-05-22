package server

import (
	"encoding/json"
	"fmt"
	"net/http"
	"quotes/models"
	"strconv"
	"strings"
	"time"
)

func (s *httpServer) PostQuote(w http.ResponseWriter, r *http.Request) {

	fmt.Printf("%v: Пришел запрос %s\n", time.Now().Format("2006-01-02 15:04:05"), r.Method)

	if r.Method != http.MethodPost {
		http.Error(w, "Неправильный метод!", http.StatusBadRequest)
		return
	}

	var quote models.Quote

	decoder := json.NewDecoder(r.Body)

	if err := decoder.Decode(&quote); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	if quote.Quote == "" || quote.Author == "" {
		http.Error(w, "Поля не должны быть пустыи", http.StatusBadRequest)
		return
	}

	if err := s.Repo.Create(r.Context(), &quote); err != nil {
		http.Error(w, "Невозможно создать цитату: "+err.Error(), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusCreated)
	w.Write([]byte("Цитата успешно создана!\n" + "Автор: " + quote.Author + "\nЦитата: " + quote.Quote))

}

func (s *httpServer) GetQuotes(w http.ResponseWriter, r *http.Request) {

	fmt.Printf("%v: Пришел запрос %s\n", time.Now().Format("2006-01-02 15:04:05"), r.Method)

	if r.Method != http.MethodGet {
		http.Error(w, "Метод должен быть GET", http.StatusMethodNotAllowed)
		return
	}

	quotes, err := s.Repo.GetAll(r.Context())
	if err != nil {
		http.Error(w, "Не удалось получить цитаты: "+err.Error(), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(quotes)

}

func (s *httpServer) GetRandomQuote(w http.ResponseWriter, r *http.Request) {

	fmt.Printf("%v: Пришел запрос %s\n", time.Now().Format("2006-01-02 15:04:05"), r.Method)

	if r.Method != http.MethodGet {
		http.Error(w, "Метод должен быть GET", http.StatusMethodNotAllowed)
		return
	}

	quote, err := s.Repo.GetRandom(r.Context())
	if err != nil {
		http.Error(w, "Не удалось получить цитату: "+err.Error(), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(quote)

}

func (s *httpServer) GetQuoteByAuthor(w http.ResponseWriter, r *http.Request) {

	fmt.Printf("%v: Пришел запрос %s\n", time.Now().Format("2006-01-02 15:04:05"), r.Method)

	if r.Method != http.MethodGet {
		http.Error(w, "Метод должен быть GET", http.StatusMethodNotAllowed)
		return
	}

	author := r.URL.Query().Get("author")
	if author == "" {
		http.Error(w, "Обязательно указать параметр автора", http.StatusBadRequest)
		return
	}

	quotes, err := s.Repo.GetByAuthor(r.Context(), author)
	if err != nil {
		http.Error(w, "Ошибка при получении цитат автора: "+err.Error(), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(quotes)
}

func (s *httpServer) DeleteQuote(w http.ResponseWriter, r *http.Request) {

	fmt.Printf("%v: Пришел запрос %s\n", time.Now().Format("2006-01-02 15:04:05"), r.Method)

	if r.Method != http.MethodDelete {
		http.Error(w, "Метод должен быть DELETE", http.StatusMethodNotAllowed)
		return
	}

	parts := strings.Split(r.URL.Path, "/")
	if len(parts) < 3 {
		http.Error(w, "Неправильный айди цитаты", http.StatusBadRequest)
		return
	}

	id, err := strconv.Atoi(parts[2])
	if err != nil {
		http.Error(w, "Неправильный формат айди цитаты", http.StatusBadRequest)
		return
	}

	if err := s.Repo.Delete(r.Context(), id); err != nil {
		http.Error(w, "Не удалось удалить цитату: "+err.Error(), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusNoContent)
	w.Write([]byte("Цитата успешно удалена!"))
}
