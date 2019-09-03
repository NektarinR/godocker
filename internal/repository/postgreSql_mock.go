package repository

import (
	"context"
	"errors"
	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/postgres"
	"time"
)

type PostgreMock struct {
	pool    []User
	logFunc FuncLogging
}

func (p *PostgreMock) InsertUser(ctx context.Context, user *User) error {
	p.pool = append(p.pool, *user)
	return nil
}

func (p *PostgreMock) GetUserById(ctx context.Context, id int) (*User, error) {
	for _, v := range p.pool {
		time.Sleep(400 * time.Millisecond)
		if v.Id == id {
			return &v, nil
		}
	}
	return nil, gorm.ErrRecordNotFound
}

func (p *PostgreMock) Fetch(ctx context.Context, offset, limit int) ([]User, error) {
	if limit > 3 {
		limit = 3
	}
	if offset > 3 {
		return nil, errors.New("data is empty")
	}
	return p.pool[offset:limit], nil
}

func (p *PostgreMock) Ping(ctx context.Context) error {
	return nil
}

func NewPostgresDBMock() (IRepository, error) {
	test := []User{
		{PrivateUser{
			Id:       1,
			CreateOn: time.Unix(10, 10)},
			PublicUser{Name: "Vasy"},
		},
		{PrivateUser{
			Id:       2,
			CreateOn: time.Unix(10, 10)},
			PublicUser{Name: "VasyVasy"},
		},
		{PrivateUser{
			Id:       3,
			CreateOn: time.Unix(10, 10)},
			PublicUser{Name: "Pety"},
		},
		{PrivateUser{
			Id:       4,
			CreateOn: time.Unix(10, 10)},
			PublicUser{Name: "PetyPety"},
		},
		{PrivateUser{
			Id:       5,
			CreateOn: time.Unix(10, 10)},
			PublicUser{Name: "Sany"},
		},
	}
	return &PostgreMock{pool: test}, nil
}
