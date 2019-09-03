package repository

import (
	"context"
	"fmt"
	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/postgres"
	uuid "github.com/satori/go.uuid"
)

type FuncLogging func(text string)

type DbConfig struct {
	Port     int
	Host     string
	DbName   string
	User     string
	Password string
}

type PostgreSql struct {
	pool    *gorm.DB
	logFunc FuncLogging
}

func (p *PostgreSql) InsertUser(ctx context.Context, user *User) error {
	u1, ok := ctx.Value("LogID").(uuid.UUID)
	if !ok {
		u1 = uuid.NewV4()
	}
	tx := p.pool.Begin()
	p.logFunc("begin tx " + u1.String())
	defer func(tx *gorm.DB) {
		if r := recover(); r != nil {
			tx.Rollback()
			p.logFunc("tx " + u1.String() + " rollback")
		}
	}(tx)
	if err := tx.Create(user).Error; err != nil {
		tx.Rollback()
		p.logFunc("tx " + u1.String() + " rollback")
		return err
	}
	tx.Commit()
	p.logFunc("tx" + u1.String() + " commit")
	return nil
}

func (p *PostgreSql) GetUserById(ctx context.Context, id int) (*User, error) {
	result := User{}
	db := p.pool.First(&result, "id = ?", id)
	if err := db.Error; err != nil {
		return nil, err
	}
	return &result, nil
}

func (p *PostgreSql) Fetch(ctx context.Context, offset, limit int) ([]User, error) {
	result := make([]User, 0, limit)
	if err := p.pool.Limit(limit).Offset(offset).Find(&result).Error; err != nil {
		return nil, err
	}
	return result, nil
}

func (p *PostgreSql) Ping(ctx context.Context) error {
	return p.pool.DB().PingContext(ctx)
}

func NewPostgreDB(config *DbConfig, fn FuncLogging) (IRepository, error) {
	conn := fmt.Sprintf("host=%s port=%d user=%s dbname=%s password=%s",
		config.Host, config.Port, config.User, config.DbName, config.Password)
	poolConn, err := gorm.Open("postgres", conn)
	if err != nil {
		return nil, err
	}
	return &PostgreSql{pool: poolConn}, nil
}
