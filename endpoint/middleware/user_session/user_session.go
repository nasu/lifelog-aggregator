package user_session

import (
	"net/http"
	"net/url"

	"github.com/labstack/echo-contrib/session"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	"golang.org/x/net/context"

	"github.com/nasu/lifelog-aggregator/constant"
	dsess "github.com/nasu/lifelog-aggregator/domain/session"
	"github.com/nasu/lifelog-aggregator/infrastructure/dynamodb"
)

type (
	Config struct {
		Skipper middleware.Skipper
		DB      *dynamodb.DB
	}
)

const (
	key = "_user_session"
)

var (
	DefaultConfig = Config{
		Skipper: func(c echo.Context) bool {
			for _, path := range constant.SKIP_USER_SESSION_MIDDLEWARE {
				if path == c.Path() {
					return true
				}
			}
			return false
		},
	}
)

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
		panic("echo: user_session middleware requires db")
	}

	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			if config.Skipper(c) {
				return next(c)
			}

			ctx := context.Background()
			sess, _ := session.Get(constant.SESSION_AUTH, c)
			sessID := sess.Values[constant.SESSION_AUTH_CONTENT_SESS_ID]
			if sessID == nil {
				return c.Redirect(http.StatusFound, constant.URL_AUTH_GOOGLE+"?redirect_uri="+url.QueryEscape(c.Request().RequestURI))
			}
			sessEntity, err := dsess.NewRepository(config.DB).Get(ctx, sessID.(string))
			if err != nil {
				return c.String(http.StatusInternalServerError, "cannot session entity")
			}
			if sessEntity.SessionID != sessID {
				return c.Redirect(http.StatusFound, constant.URL_AUTH_GOOGLE+"?redirect_uri="+url.QueryEscape(c.Request().RequestURI))
			}
			c.Set(key, sessEntity)
			return next(c)
		}
	}
}
