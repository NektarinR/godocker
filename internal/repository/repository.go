package repository

import (
	"context"
	"time"
)

type User struct {
	PrivateUser
	PublicUser
}

type PrivateUser struct {
	Id       int       `gorm:"column:id" json:"id, int" reform:"id,pk"`
	CreateOn time.Time `gorm:"column:created_on" json:"create_on" reform:"created_on"`
}

type PublicUser struct {
	Name string `gorm:"column:name" json:"name" reform:"name"`
}

type IRepository interface {
	InsertUser(ctx context.Context, user *User) error
	GetUserById(ctx context.Context, id int) (*User, error)
	Fetch(ctx context.Context, offset, limit int) ([]User, error)
	Ping(ctx context.Context) error
}
