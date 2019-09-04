package repository

import (
	"context"
	"github.com/DATA-DOG/go-sqlmock"
	"github.com/jinzhu/gorm"
	uuid "github.com/satori/go.uuid"
	"log"
	"reflect"
	"regexp"
	"testing"
	"time"
)

type TestDB struct {
	mock sqlmock.Sqlmock
	repo IRepository
}

var (
	p        = TestDB{}
	testuser = []User{
		{PrivateUser{
			Id:       1,
			CreateOn: time.Now()},
			PublicUser{Name: "Vasy"},
		},
		{PrivateUser{
			Id:       2,
			CreateOn: time.Now()},
			PublicUser{Name: "VasyVasy"},
		},

		{PrivateUser{
			Id:       3,
			CreateOn: time.Now()},
			PublicUser{Name: "VasyVasyVasy"},
		},
	}
)

func Setup() {
	db, mock, _ := sqlmock.New()
	p.mock = mock
	testDb, _ := gorm.Open("postgres", db)
	testDb.LogMode(false)
	repo := &PostgreSql{pool: testDb, logFunc: func(text string) {
		log.Println(text)
	}}
	p.repo = repo
}

func TestPostgreSql_Fetch(t *testing.T) {
	Setup()
	p.mock.ExpectQuery(regexp.QuoteMeta(`SELECT * FROM "users" LIMIT 3 OFFSET 0 `)).
		WillReturnRows(sqlmock.NewRows([]string{"id", "created_on", "name"}).
			AddRow(testuser[0].Id, testuser[0].CreateOn, testuser[0].Name).
			AddRow(testuser[1].Id, testuser[1].CreateOn, testuser[1].Name).
			AddRow(testuser[2].Id, testuser[2].CreateOn, testuser[2].Name))
	res, err := p.repo.Fetch(context.Background(), 0, 3)
	if err != nil {
		t.Errorf("expected nil got %s", err)
	}
	if !reflect.DeepEqual(testuser, res) {
		t.Errorf("expected %v got %v", testuser, res)
	}
}

func TestPostgreSql_Fetch_1(t *testing.T) {
	Setup()
	p.mock.ExpectQuery(regexp.QuoteMeta(`SELECT * FROM "users" LIMIT 2 OFFSET 0`)).
		WillReturnRows(sqlmock.NewRows([]string{"id", "created_on", "name"}).
			AddRow(testuser[0].Id, testuser[0].CreateOn, testuser[0].Name).
			AddRow(testuser[1].Id, testuser[1].CreateOn, testuser[1].Name))
	res, err := p.repo.Fetch(context.Background(), 0, 2)
	if err != nil {
		t.Errorf("expected nil got %s", err)
	}
	if !reflect.DeepEqual(res, testuser[:2]) {
		t.Errorf("expected %v \ngot %v", testuser[:2], res)
	}
}

//return 2
func TestPostgreSql_Fetch2(t *testing.T) {
	Setup()
	p.mock.ExpectQuery(regexp.QuoteMeta(`SELECT * FROM "users" LIMIT 2 OFFSET 1`)).
		WillReturnRows(sqlmock.NewRows([]string{"id", "created_on", "name"}).
			AddRow(testuser[1].Id, testuser[1].CreateOn, testuser[1].Name).
			AddRow(testuser[2].Id, testuser[2].CreateOn, testuser[2].Name))
	res, err := p.repo.Fetch(context.Background(), 1, 2)
	if err != nil {
		t.Errorf("expected nil got %s", err)
	}
	if !reflect.DeepEqual(res, testuser[1:3]) {
		t.Errorf("expected %v \ngot %v", testuser[1:3], res)
	}
}

//return 1
func TestPostgreSql_Fetch3(t *testing.T) {
	Setup()
	p.mock.ExpectQuery(regexp.QuoteMeta(`SELECT * FROM "users" LIMIT 3 OFFSET 2`)).
		WillReturnRows(sqlmock.NewRows([]string{"id", "created_on", "name"}).
			AddRow(testuser[2].Id, testuser[2].CreateOn, testuser[2].Name))
	res, err := p.repo.Fetch(context.Background(), 2, 3)
	if err != nil {
		t.Errorf("expected nil got %s", err)
	}
	if !reflect.DeepEqual(res, testuser[2:3]) {
		t.Errorf("expected %v \ngot %v", testuser[2:3], res)
	}
}

//return 0
func TestPostgreSql_Fetch4(t *testing.T) {
	Setup()
	p.mock.ExpectQuery(regexp.QuoteMeta(`SELECT * FROM "users" LIMIT 10 OFFSET 4`)).
		WillReturnRows(sqlmock.NewRows([]string{"id", "created_on", "name"}))
	res, err := p.repo.Fetch(context.Background(), 4, 10)
	if err != nil {
		t.Errorf("expected nil got %s", err)
	}
	if !reflect.DeepEqual(res, []User{}) {
		t.Errorf("expected %v \ngot %v", []User{}, res)
	}
}

func TestPostgreSql_FetchError(t *testing.T) {
	Setup()
	p.mock.ExpectQuery(regexp.QuoteMeta(`SELECT * FROM "users" LIMIT 3 OFFSET 0`)).
		WillReturnError(gorm.ErrRecordNotFound)
	res, err := p.repo.Fetch(context.Background(), 0, 3)
	if err == nil {
		t.Errorf("expected %v got %s", res, err)
	}
	if reflect.DeepEqual(res, []User{}) {
		t.Errorf("expected %v \ngot %v", []User{}, res)
	}
}

func TestPostgreSql_GetUserById(t *testing.T) {
	Setup()
	p.mock.ExpectQuery(regexp.QuoteMeta(`SELECT * FROM "users" WHERE (id = $1)`)).
		WithArgs(2).
		WillReturnRows(sqlmock.NewRows([]string{"id", "created_on", "name"}).
			AddRow(testuser[1].Id, testuser[1].CreateOn, testuser[1].Name))
	res, err := p.repo.GetUserById(context.Background(), 2)
	if err != nil {
		t.Errorf("expected nil got %v", err)
	}
	if !reflect.DeepEqual(*res, testuser[1]) {
		t.Errorf("expected %v \ngot %v", testuser[1], res)
	}
}

func TestPostgreSql_GetUserById_Error(t *testing.T) {
	Setup()
	p.mock.ExpectQuery(regexp.QuoteMeta(`SELECT * FROM "users" WHERE (id = $1)`)).
		WithArgs(2).
		WillReturnRows(sqlmock.NewRows([]string{"id", "created_on", "name"}))
	res, err := p.repo.GetUserById(context.Background(), 2)
	if err != gorm.ErrRecordNotFound {
		t.Errorf("expected %v got %v", gorm.ErrRecordNotFound, err)
	}
	if res != nil {
		t.Errorf("expected %v \ngot %v", nil, res)
	}
}

func TestPostgreSql_Ping(t *testing.T) {
	Setup()
	p.mock.ExpectQuery(regexp.QuoteMeta(`SELECT * FROM "users" WHERE (id = $1)`)).
		WithArgs(2).
		WillReturnRows(sqlmock.NewRows([]string{"id", "created_on", "name"}))
	err := p.repo.Ping(context.Background())
	if err != nil {
		t.Errorf("expected %v got %s", gorm.ErrRecordNotFound, err)
	}
}

func TestPostgreSql_InsertUser_Success(t *testing.T) {
	Setup()
	time := time.Now()
	user := User{}
	user.Name = "Vasy"
	user.CreateOn = time
	strQuery := regexp.QuoteMeta(`INSERT  INTO "users" ("created_on","name") 
		VALUES ($1,$2) RETURNING "users"."id"`)
	p.mock.ExpectBegin()
	p.mock.ExpectQuery(strQuery).
		WithArgs(user.CreateOn, user.Name).
		WillReturnRows(sqlmock.NewRows([]string{"id"}).
			AddRow(1))
	p.mock.ExpectCommit()
	ctx := context.WithValue(context.Background(), "LogID", uuid.NewV4())
	err := p.repo.InsertUser(ctx, &user)
	if err != nil {
		t.Errorf("expected nil got\n %s", err)
	}
}

func TestPostgreSql_InsertUser_Err(t *testing.T) {
	Setup()
	time := time.Now()
	user := User{}
	user.Name = "Vasy"
	user.CreateOn = time
	strQuery := regexp.QuoteMeta(`INSERT  INTO "users" ("created_on","name") 
		VALUES ($1,$2) RETURNING "users"."id"`)
	p.mock.ExpectBegin()
	p.mock.ExpectQuery(strQuery).
		WithArgs(user.CreateOn, user.Name).
		WillReturnError(gorm.ErrInvalidTransaction)
	p.mock.ExpectRollback()
	ctx := context.WithValue(context.Background(), "LogID", uuid.NewV4())
	_ = p.repo.InsertUser(ctx, &user)
	if err := p.mock.ExpectationsWereMet(); err != nil {
		t.Errorf("expected nil, got:\n %s", err)
	}
}

func TestPostgreSql_InsertUser_Err_Uuid(t *testing.T) {
	Setup()
	time := time.Now()
	user := User{}
	user.Name = "Vasy"
	user.CreateOn = time
	strQuery := regexp.QuoteMeta(`INSERT  INTO "users" ("created_on","name") 
		VALUES ($1,$2) RETURNING "users"."id"`)
	p.mock.ExpectBegin()
	p.mock.ExpectQuery(strQuery).
		WithArgs(user.CreateOn, user.Name).
		WillReturnError(gorm.ErrInvalidTransaction)
	p.mock.ExpectRollback()
	ctx := context.WithValue(context.Background(), "Lg", uuid.NewV4())
	_ = p.repo.InsertUser(ctx, &user)
	if err := p.mock.ExpectationsWereMet(); err != nil {
		t.Errorf("expected nil, got:\n %s", err)
	}
}

//func TestNewPostgreDB(t *testing.T) {
//	dbConf := DbConfig{
//		Port:     5433,
//		Host:     "localhost",
//		DbName:   "test",
//		User:     "postgres",
//		Password: "12345",
//	}
//	_, err := NewPostgreDB(&dbConf, func(text string) {
//		log.Println(text)
//	})
//	if err != nil {
//		t.Errorf("expected nil, got:\n %s", err)
//	}
//}

func TestNewPostgreDB_Error(t *testing.T) {
	dbConf := DbConfig{
		Port:     5432,
		Host:     "db",
		DbName:   "test",
		User:     "postgres",
		Password: "",
	}
	_, err := NewPostgreDB(&dbConf, func(text string) {
		log.Println(text)
	})

	if err == nil {
		t.Errorf("expected nil, got:\n %#v", err)
	}
}
