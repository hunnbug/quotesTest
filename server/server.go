package server

import (
	"fmt"
	"net/http"
	"quotes/models"

	"github.com/gorilla/mux"
)

type Server interface {
	Start() error
	PostQuote(w http.ResponseWriter, r *http.Request)
	GetQuotes(w http.ResponseWriter, r *http.Request)
	GetRandomQuote(w http.ResponseWriter, r *http.Request)
	GetQuoteByAuthor(w http.ResponseWriter, r *http.Request)
	DeleteQuote(w http.ResponseWriter, r *http.Request)
}

type httpServer struct {
	Addr string
	Repo models.QuoteRepo
}

func NewHttpServer(addr string, repo models.QuoteRepo) Server {
	return &httpServer{Addr: addr, Repo: repo}
}

func (s *httpServer) Start() error {

	r := mux.NewRouter()

	r.HandleFunc("/quotes", s.PostQuote).Methods(http.MethodPost)
	r.HandleFunc("/quotes", s.GetQuoteByAuthor).Queries("author", "{author}").Methods(http.MethodGet)
	r.HandleFunc("/quotes", s.GetQuotes).Methods(http.MethodGet)
	r.HandleFunc("/quotes/random", s.GetRandomQuote).Methods(http.MethodGet)
	r.HandleFunc("/quotes/{id}", s.DeleteQuote).Methods(http.MethodDelete)

	fmt.Printf("Сервер запущен на порте %s\n", s.Addr)

	err := http.ListenAndServe(s.Addr, r)

	return err
}
