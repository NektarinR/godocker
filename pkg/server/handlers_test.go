package server

import (
	"encoding/json"
	"github.com/NektarinR/godocker/internal/repository"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"
)

type TestCase struct {
	Method         string
	Url            string
	RequestBody    string
	ResponseStatus int
	ResponseBody   string
}

func TestServer_HandlePing_200(t *testing.T) {
	srv := Server{}
	srv.InitRouters()
	testCase := &TestCase{
		Method:         "GET",
		Url:            "localhost:8081/ping",
		ResponseStatus: 200,
		ResponseBody:   "",
	}
	w := httptest.NewRecorder()
	req, _ := http.NewRequest(testCase.Method, testCase.Url, nil)
	srv.mux.ServeHTTP(w, req)
	if w.Code != testCase.ResponseStatus {
		t.Errorf("wrong responce code, got %d expected %d\n",
			w.Code, testCase.ResponseStatus)
	}
}

func TestServer_HandlePing_405(t *testing.T) {
	srv := Server{}
	srv.InitRouters()
	testCase := &TestCase{
		Method:         "HEAD",
		Url:            "http://localhost:8081/ping",
		ResponseStatus: http.StatusMethodNotAllowed,
		ResponseBody:   "",
	}
	w := httptest.NewRecorder()
	req, _ := http.NewRequest(testCase.Method, testCase.Url, nil)
	srv.mux.ServeHTTP(w, req)
	if w.Code != testCase.ResponseStatus {
		t.Errorf("wrong responce code, got %d expected %d\n",
			w.Code, testCase.ResponseStatus)
	}
}

func TestServer_HandleGetUserById_Success(t *testing.T) {
	srv := Server{}
	srv.InitRouters()
	srv.db, _ = repository.NewPostgresDBMock()
	tms := time.Unix(10, 10)
	userTest, _ := json.Marshal(repository.User{
		PrivateUser: repository.PrivateUser{Id: 1, CreateOn: tms},
		PublicUser:  repository.PublicUser{Name: "Vasy"},
	})
	testCase := &TestCase{
		Method:         "GET",
		Url:            "http://localhost:8081/users/1",
		ResponseStatus: 200,
		ResponseBody:   string(userTest),
	}
	w := httptest.NewRecorder()
	req, _ := http.NewRequest(testCase.Method, testCase.Url, nil)
	srv.mux.ServeHTTP(w, req)
	if w.Code != testCase.ResponseStatus {
		t.Errorf("wrong responce code, got %d expected %d\n",
			w.Code, testCase.ResponseStatus)
	}
	body, err := ioutil.ReadAll(w.Body)
	if err != nil {
		t.Errorf("expected nil, got %v\n", err)
	}
	if string(body) != testCase.ResponseBody {
		t.Errorf("expected %v, got %v\n", testCase.ResponseBody, string(body))
	}
}

func TestServer_HandleGetUserById_ErrorBadId(t *testing.T) {
	srv := Server{}
	srv.InitRouters()
	srv.db, _ = repository.NewPostgresDBMock()
	testCase := &TestCase{
		Method:         "GET",
		Url:            "http://localhost:8081/users/999999999999999999999999",
		ResponseStatus: http.StatusBadRequest,
		ResponseBody:   "bad id\n",
	}
	w := httptest.NewRecorder()
	req, _ := http.NewRequest(testCase.Method, testCase.Url, nil)
	srv.mux.ServeHTTP(w, req)
	if w.Code != testCase.ResponseStatus {
		t.Errorf("wrong responce code, got %d expected %d\n",
			w.Code, testCase.ResponseStatus)
	}
	body, err := ioutil.ReadAll(w.Body)
	if err != nil {
		t.Errorf("expected nil, got %v\n", err)
	}
	if string(body) != testCase.ResponseBody {
		t.Errorf("expected %v, got %v\n", testCase.ResponseBody, string(body))
	}
}

func TestServer_HandleGetUserById_ErrorServerBusy(t *testing.T) {
	srv := Server{}
	srv.InitRouters()
	srv.db, _ = repository.NewPostgresDBMock()
	testCase := &TestCase{
		Method:         "GET",
		Url:            "http://localhost:8081/users/5",
		ResponseStatus: http.StatusInternalServerError,
		ResponseBody:   "server is busy\n",
	}
	w := httptest.NewRecorder()
	req, _ := http.NewRequest(testCase.Method, testCase.Url, nil)
	srv.mux.ServeHTTP(w, req)
	if w.Code != testCase.ResponseStatus {
		t.Errorf("wrong responce code, got %d expected %d\n",
			w.Code, testCase.ResponseStatus)
	}
	body, err := ioutil.ReadAll(w.Body)
	if err != nil {
		t.Errorf("expected nil, got %v\n", err)
	}
	if string(body) != testCase.ResponseBody {
		t.Errorf("expected %v, got %v\n", testCase.ResponseBody, string(body))
	}
}

func TestServer_HandleGetUserById_Error(t *testing.T) {
	srv := Server{}
	srv.InitRouters()
	srv.db, _ = repository.NewPostgresDBMock()
	testCase := &TestCase{
		Method:         "GET",
		Url:            "http://localhost:8081/users/3",
		ResponseStatus: http.StatusInternalServerError,
		ResponseBody:   "server is busy\n",
	}
	w := httptest.NewRecorder()
	req, _ := http.NewRequest(testCase.Method, testCase.Url, nil)
	srv.mux.ServeHTTP(w, req)
	if w.Code != testCase.ResponseStatus {
		t.Errorf("wrong responce code, got %d expected %d\n",
			w.Code, testCase.ResponseStatus)
	}
	body, err := ioutil.ReadAll(w.Body)
	if err != nil {
		t.Errorf("expected nil, got %v\n", err)
	}
	if string(body) != testCase.ResponseBody {
		t.Errorf("expected %v, got %v\n", testCase.ResponseBody, string(body))
	}
}
