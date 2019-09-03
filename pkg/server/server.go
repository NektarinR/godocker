package server

import (
	"context"
	"github.com/NektarinR/godocker/internal/repository"
	mx "github.com/gorilla/mux"
	"log"
	"net/http"
	"os"
	"os/signal"
	"strconv"
	"syscall"
	"time"
)

type Server struct {
	mux *mx.Router
	db  repository.IRepository
}

func Logging(text string) {
	log.Println(text)
}

func (p *Server) InitRouters() {
	p.mux = mx.NewRouter()
	p.mux.HandleFunc("/ping", p.HandlePing).
		Methods(http.MethodGet)
	p.mux.HandleFunc("/users", p.HandleGetUsers).
		Queries("offset", "{offset:[0-9]+}").
		Queries("limit", "{limit:[0-9]+}").
		Methods(http.MethodGet)
	p.mux.HandleFunc("/users/{id:[0-9]+}", p.HandleGetUserById).
		Methods(http.MethodGet)
	p.mux.HandleFunc("/users/", p.HandleInsertUser).
		Methods(http.MethodPost)
	p.mux.Use(p.loggingMiddleware)
}

func (p *Server) InitDb() {
	conf := &repository.DbConfig{
		Host:     "localhost",
		Port:     5432,
		User:     "postgres",
		Password: "",
		DbName:   "testing",
	}
	p.db, _ = repository.NewPostgreDB(conf, Logging)
}

func (p *Server) Run(port int) {
	exit := make(chan os.Signal, 1)
	signal.Notify(exit, syscall.SIGINT)

	p.InitDb()
	p.InitRouters()

	srv := &http.Server{
		Addr:         ":" + strconv.Itoa(port),
		Handler:      p.mux,
		WriteTimeout: 15 * time.Second,
		ReadTimeout:  15 * time.Second,
		IdleTimeout:  60 * time.Second,
	}

	go func(serv *http.Server, exitHttp <-chan os.Signal) {
		<-exitHttp
		log.Println("Server is closing")
		ctx, cancel := context.WithTimeout(context.Background(),
			10*time.Second)
		defer cancel()
		if err := serv.Shutdown(ctx); err != nil && err != http.ErrServerClosed {

		}
	}(srv, exit)

	log.Println("Server is started")
	if err := srv.ListenAndServe(); err != nil {

	}
	log.Println("Server is closed")
}
