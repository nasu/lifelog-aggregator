package google

import (
	"context"
	"fmt"
	"net/http"
	"os"
	"time"

	"github.com/google/uuid"
	"github.com/gorilla/sessions"
	"github.com/labstack/echo-contrib/session"
	"github.com/labstack/echo/v4"
	"golang.org/x/oauth2"
	v2 "google.golang.org/api/oauth2/v2"
	"google.golang.org/api/option"

	"github.com/nasu/lifelog-aggregator/constant"
	dsess "github.com/nasu/lifelog-aggregator/domain/session"
	"github.com/nasu/lifelog-aggregator/endpoint/middleware/database"
)

const (
	authURL  = "https://accounts.google.com/o/oauth2/v2/auth"
	tokenURL = "https://www.googleapis.com/oauth2/v4/token"
	state    = "G4=y+1bbE&@9"
)

func Index(c echo.Context) error {
	sess, _ := session.Get(constant.SESSION_FLASH, c)
	sess.Options = &sessions.Options{
		Path:     "/",
		MaxAge:   5 * 60,
		HttpOnly: true,
		//TODO
		//Secure: true,
		//SameSite: http.SameSiteStrictMode,
	}
	sess.AddFlash(c.QueryParam("redirect_uri"))
	sess.Save(c.Request(), c.Response())

	auth := googleAuth{
		clientID:     os.Getenv("GOOGLE_CLIENT_ID"),
		clientSecret: os.Getenv("GOOGLE_CLIENT_SECRET"),
	}
	config := auth.makeConfig()
	return c.Redirect(http.StatusFound, config.AuthCodeURL(state, oauth2.AccessTypeOffline, oauth2.ApprovalForce))
}

func Cb(c echo.Context) error {
	if c.QueryParam("state") != state {
		return c.String(http.StatusBadRequest, "state is wrong")
	}

	flashSess, _ := session.Get(constant.SESSION_FLASH, c)
	flashes := flashSess.Flashes()
	//TODO: 中身安全かチェック
	var redirectURI string
	if len(flashes) > 0 {
		redirectURI = constant.URL + flashes[0].(string)
	} else {
		redirectURI = constant.URL
	}
	flashSess.Save(c.Request(), c.Response())

	auth := googleAuth{
		clientID:     os.Getenv("GOOGLE_CLIENT_ID"),
		clientSecret: os.Getenv("GOOGLE_CLIENT_SECRET"),
	}
	config := auth.makeConfig()

	ctx := context.Background()
	token, err := config.Exchange(ctx, c.QueryParam("code"))
	if err != nil {
		return err
	}
	if !token.Valid() {
		return fmt.Errorf("token is invalid")
	}
	service, err := v2.NewService(ctx, option.WithTokenSource(config.TokenSource(ctx, token)))
	if err != nil {
		return err
	}
	tokenInfo, err := service.Tokeninfo().AccessToken(token.AccessToken).Context(ctx).Do()
	if err != nil {
		return err
	}

	sessID := uuid.NewString()
	db, err := database.Get(c)
	if err != nil {
		return c.String(http.StatusInternalServerError, err.Error())
	}
	err = dsess.NewRepository(db).Save(ctx, &dsess.Session{
		SessionID: sessID,
		UserID:    tokenInfo.UserId,
		Email:     tokenInfo.Email,
		CreatedAt: time.Now(),
	})
	if err != nil {
		return c.String(http.StatusInternalServerError, err.Error())
	}

	authSess, _ := session.Get(constant.SESSION_AUTH, c)
	authSess.Options = &sessions.Options{
		Path:     "/",
		MaxAge:   24 * 60 * 60,
		HttpOnly: true,
		//TODO
		//Secure: true,
		//SameSite: http.SameSiteStrictMode,
	}
	authSess.Values[constant.SESSION_AUTH_CONTENT_SESS_ID] = sessID
	authSess.Save(c.Request(), c.Response())
	return c.Redirect(http.StatusFound, redirectURI)
}

type googleAuth struct {
	clientID     string
	clientSecret string
}

func (g googleAuth) makeConfig() *oauth2.Config {
	config := &oauth2.Config{
		ClientID:     g.clientID,
		ClientSecret: g.clientSecret,
		Endpoint: oauth2.Endpoint{
			AuthURL:  authURL,
			TokenURL: tokenURL,
		},
		Scopes:      []string{"openid", "profile", "email"},
		RedirectURL: constant.URL_AUTH_GOOGLE_CALLBACK,
	}
	return config
}
