package database

import (
	"fmt"

	"github.com/gorilla/context"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"

	"github.com/nasu/lifelog-aggregator/infrastructure/dynamodb"
)

type (
	Config struct {
		Skipper middleware.Skipper
		DB      *dynamodb.DB
	}
)

const (
	key = "_dynamodb_client"
)

var (
	DefaultConfig = Config{
		Skipper: middleware.DefaultSkipper,
	}
)

func Get(c echo.Context) (*dynamodb.DB, error) {
	db := c.Get(key)
	if db == nil {
		return nil, fmt.Errorf("db not found")
	}
	return db.(*dynamodb.DB), nil
}

func Middleware(db *dynamodb.DB) echo.MiddlewareFunc {
	c := DefaultConfig
	c.DB = db
	return MiddlewareWithConfig(c)
}

func MiddlewareWithConfig(config Config) echo.MiddlewareFunc {
	if config.Skipper == nil {
		config.Skipper = DefaultConfig.Skipper
	}
	if config.DB == nil {
		panic("echo: db middleware requires db")
	}

	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			if config.Skipper(c) {
				return next(c)
			}
			defer context.Clear(c.Request())
			c.Set(key, config.DB)
			return next(c)
		}
	}
}
