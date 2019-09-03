package server

import (
	"context"
	uuid "github.com/satori/go.uuid"
	"log"
	"net/http"
)

type CustResponce struct {
	w          http.ResponseWriter
	statusCode int
}

func (p *CustResponce) Write(answer []byte) (int, error) {
	return p.w.Write(answer)
}

func (p *CustResponce) WriteHeader(statusCode int) {
	p.w.WriteHeader(statusCode)
	p.statusCode = statusCode
}

func (p *CustResponce) Header() http.Header {
	return p.w.Header()
}

func (p *Server) loggingMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		u1 := uuid.NewV4()
		ctx := context.WithValue(r.Context(),
			"LogID", u1)
		resp := &CustResponce{w: w}
		next.ServeHTTP(resp, r.WithContext(ctx))
		log.Printf("%s %s %s %s %d\n", u1.String(), r.RemoteAddr,
			r.Method, r.RequestURI, resp.statusCode)
	})
}
