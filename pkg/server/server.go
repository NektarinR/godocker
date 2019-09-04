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
	log.Println("Старт инициализация routes")
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
	log.Println("Конец инициализации routes")
}

func (p *Server) InitDb() (err error) {
	log.Println("Старт инициализация соединения с db")
	conf := &repository.DbConfig{
		Host:     "db",
		Port:     5432,
		User:     "postgres",
		Password: "12345",
		DbName:   "test",
	}
	p.db, err = repository.NewPostgreDB(conf, Logging)
	if err != nil {
		log.Printf("Ошибка при соединение с БД %v\n", err)
	}
	log.Println("Конец инициализация соединения с db")
	return nil
}

func (p *Server) Run(port int) {
	log.Printf("Запуск http сервера на порту %d\n", port)
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
		log.Println("Сервер останавливается...")
		ctx, cancel := context.WithTimeout(context.Background(),
			10*time.Second)
		defer cancel()
		if err := serv.Shutdown(ctx); err != nil && err != http.ErrServerClosed {

		}
	}(srv, exit)

	log.Println("Сервер запущен")
	if err := srv.ListenAndServe(); err != nil {

	}
	log.Println("Сервер остановлен")
}
