package main

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"os"
	"time"

	"github.com/labstack/echo"
	"github.com/labstack/echo/middleware"
	"golang.org/x/oauth2"
	v2 "google.golang.org/api/oauth2/v2"

	"github.com/nasu/lifelog-aggregator/oura"
	"github.com/nasu/lifelog-aggregator/toggl"
	"github.com/nasu/lifelog-aggregator/util/logger"
)

var logg *logger.Logger

func init() {
	logLevel := os.Getenv("LOGLEVEL")
	logg = logger.NewLoggerWithStringLogLevel(logLevel)
}

func main() {
	e := echo.New()
	e.Use(middleware.Logger())
	e.Use(middleware.Recover())
	e.GET("/diary", diary)
	e.GET("/auth/cb", authCb)
	e.Logger.Fatal(e.Start(":8080"))
}

func diary(c echo.Context) error {
	auth := c.Request().Header.Get(echo.HeaderAuthorization)
	if auth == "" {
		config := &oauth2.Config{
			ClientID:     os.Getenv("GOOGLE_CLIENT_ID"),
			ClientSecret: os.Getenv("GOOGLE_CLIENT_SECRET"),
			Endpoint: oauth2.Endpoint{
				AuthURL:  "https://accounts.google.com/o/oauth2/v2/auth",
				TokenURL: "https://www.googleapis.com/oauth2/v4/token",
			},
			Scopes:      []string{"openid", "profile", "email"},
			RedirectURL: "http://localhost:8080/auth/cb",
		}
		state := "abcdefg"
		return c.Redirect(http.StatusFound, config.AuthCodeURL(state, oauth2.AccessTypeOffline, oauth2.ApprovalForce))
	}

	now := time.Now()
	total := make(map[string]map[string]int)
	for i := 1; i <= 7; i++ {
		date := now.Add(time.Hour * 24 * -1 * time.Duration(i)).Format("2006-01-02")
		projectTimes := fromToggle(date)
		projectTimes["睡眠"] += fromOura(date)
		total[date] = projectTimes
	}
	return c.String(http.StatusOK, fmt.Sprintf("%#v", total))
}

func authCb(c echo.Context) error {
	config := &oauth2.Config{
		ClientID:     os.Getenv("GOOGLE_CLIENT_ID"),
		ClientSecret: os.Getenv("GOOGLE_CLIENT_SECRET"),
		Endpoint: oauth2.Endpoint{
			AuthURL:  "https://accounts.google.com/o/oauth2/v2/auth",
			TokenURL: "https://www.googleapis.com/oauth2/v4/token",
		},
		Scopes:      []string{"openid", "profile"},
		RedirectURL: "http://localhost:8080/auth/cb",
	}
	ctx := context.Background()
	token, err := config.Exchange(ctx, c.QueryParam("code"))
	if err != nil {
		return err
	}
	if !token.Valid() {
		return fmt.Errorf("token is invalid")
	}
	service, err := v2.New(config.Client(ctx, token))
	if err != nil {
		return err
	}
	tokenInfo, err := service.Tokeninfo().AccessToken(token.AccessToken).Context(ctx).Do()
	if err != nil {
		return err
	}
	return c.String(http.StatusOK, fmt.Sprintf("%#v\n\n%s\n%s\n%s\n%s", tokenInfo, token.AccessToken, token.RefreshToken, token.Expiry, token.TokenType))
}

func fromToggle(date string) map[string]int {
	c := &toggl.Client{
		UserAgent:   os.Getenv("TOGGL_USER_AGENT"),
		WorkSpaceID: os.Getenv("TOGGL_WORKSPACE_ID"),
		ApiToken:    os.Getenv("TOGGL_API_TOKEN"),
		Logger:      logg,
	}
	details, err := c.GetDetails(date, date)
	if err != nil {
		log.Fatal(err)
	}

	projectTimes := make(map[string]int)
	for _, d := range details {
		// d.Duration is in milliseconds so that converts to seconds
		projectTimes[d.Project] += d.Duration / 1000
	}
	return projectTimes
}

func fromOura(date string) int {
	c := &oura.Client{
		AccessToken: os.Getenv("OURA_ACCESS_TOKEN"),
		Logger:      logg,
	}
	sleep, err := c.SleepOneDay(date)
	if err != nil {
		log.Fatal(err)
	}
	return sleep.Duration
}
