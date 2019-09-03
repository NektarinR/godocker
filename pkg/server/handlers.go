package server

import (
	"context"
	"encoding/json"
	"errors"
	"github.com/NektarinR/godocker/internal/repository"
	mx "github.com/gorilla/mux"
	"net/http"
	"strconv"
	"strings"
	"time"
)

func (p *Server) HandlePing(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusOK)
}

//Get method - /users?offset=12&limit=12
func (p *Server) HandleGetUsers(w http.ResponseWriter, r *http.Request) {
	offset, limit, err := parseURL(r)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	ctx, _ := context.WithTimeout(r.Context(), 2*time.Second)
	usrRes := make(chan []repository.User, 1)
	go func(insideCtx context.Context, res chan<- []repository.User) {
		usrs, err := p.db.Fetch(insideCtx, offset, limit)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		res <- usrs
	}(ctx, usrRes)
	select {
	case users := <-usrRes:
		result, err := json.Marshal(users)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		w.Write(result)
	case <-ctx.Done():
		http.Error(w, "server is busy", http.StatusInternalServerError)
		return
	}
}

//POST method - /users/
func (p *Server) HandleInsertUser(w http.ResponseWriter, r *http.Request) {
	decoder := json.NewDecoder(r.Body)
	decoder.DisallowUnknownFields()
	var pubUsr repository.PublicUser
	err := decoder.Decode(&pubUsr)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	usr := &repository.User{
		PublicUser: pubUsr,
	}
	tmpUsr := make(chan struct{}, 1)
	ctx, _ := context.WithTimeout(r.Context(), 2*time.Second)
	go func(insideCtx context.Context, res chan<- struct{}) {
		err := p.db.InsertUser(insideCtx, usr)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		res <- struct{}{}
	}(ctx, tmpUsr)
	select {
	case <-tmpUsr:
		w.WriteHeader(http.StatusOK)
	case <-ctx.Done():
		http.Error(w, "server is busy", http.StatusInternalServerError)
		return
	}
}

//get method - /users/{id}
func (p *Server) HandleGetUserById(w http.ResponseWriter, r *http.Request) {
	vars := mx.Vars(r)
	id, err := strconv.Atoi(vars["id"])
	if err != nil {
		http.Error(w, "bad id", http.StatusBadRequest)
		return
	}
	ctx, _ := context.WithTimeout(r.Context(), 2*time.Second)
	userChan := make(chan *repository.User, 1)
	exitRequest := make(chan struct{}, 1)
	go func(insideCtx context.Context,
		res chan<- *repository.User, exit chan<- struct{}) {
		usr, err := p.db.GetUserById(ctx, id)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			exit <- struct{}{}
		}
		res <- usr
	}(ctx, userChan, exitRequest)
	select {
	case <-exitRequest:
		return
	case usr := <-userChan:
		encode, err := json.Marshal(usr)
		if err != nil {
			http.Error(w, "can't json", http.StatusInternalServerError)
			return
		}
		w.WriteHeader(http.StatusOK)
		w.Write(encode)
	case <-ctx.Done():
		http.Error(w, "server is busy", http.StatusInternalServerError)
		return
	}
}

func parseURL(r *http.Request) (int, int, error) {
	vars := mx.Vars(r)
	offsetSTR, ok := vars["offset"]
	if !ok {
		return 0, 0, errors.New("bad query")
	}
	offsetSTR = strings.TrimSpace(offsetSTR)
	var offset int = 0
	if offsetSTR != "" {
		tmp, err := strconv.Atoi(offsetSTR)
		if err != nil {
			return 0, 0, errors.New("bad offset")
		}
		offset = tmp
	}
	limitStr, ok := vars["limit"]
	if !ok {
		return 0, 0, errors.New("bad query")
	}
	limitStr = strings.TrimSpace(limitStr)
	var limit int = 0
	if limitStr != "" {
		tmp, err := strconv.Atoi(limitStr)
		if err != nil {
			return 0, 0, errors.New("bad limit")
		}
		if tmp > 25 {
			tmp = 25
		}
		limit = tmp
	}
	return offset, limit, nil
}
