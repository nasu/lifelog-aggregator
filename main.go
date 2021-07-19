package main

import (
	"fmt"
	"html/template"
	"io"
	"log"
	"net/http"
	"os"
	"time"

	"github.com/gorilla/sessions"
	"github.com/labstack/echo-contrib/session"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	"golang.org/x/net/context"

	"github.com/nasu/lifelog-aggregator/constant"
	"github.com/nasu/lifelog-aggregator/domain/oura"
	"github.com/nasu/lifelog-aggregator/domain/toggl"
	"github.com/nasu/lifelog-aggregator/endpoint/auth/google"
	"github.com/nasu/lifelog-aggregator/endpoint/middleware/database"
	"github.com/nasu/lifelog-aggregator/endpoint/middleware/user_session"
	"github.com/nasu/lifelog-aggregator/infrastructure/dynamodb"
	"github.com/nasu/lifelog-aggregator/util/logger"
)

var logg *logger.Logger

type TemplateRenderer struct {
	templates *template.Template
}

func (t *TemplateRenderer) Render(w io.Writer, name string, data interface{}, c echo.Context) error {
	return t.templates.ExecuteTemplate(w, name, data)
}

func init() {
	logLevel := os.Getenv("LOGLEVEL")
	logg = logger.NewLoggerWithStringLogLevel(logLevel)

	if os.Getenv("DYNAMODB_URL") == "" {
		panic("DYNAMODB_URL is required")
	}
	if os.Getenv("DYNAMODB_REGION") == "" {
		panic("DYNAMODB_REGION is required")
	}
	if os.Getenv("SESSION_STORE_KEY") == "" {
		panic("SESSION_STORE_KEY is required")
	}

	if os.Getenv("GOOGLE_CLIENT_ID") == "" {
		panic("GOOGLE_CLIENT_ID is required")
	}
	if os.Getenv("GOOGLE_CLIENT_SECRET") == "" {
		panic("GOOGLE_CLIENT_SECRET is required")
	}

	if os.Getenv("TOGGL_USER_AGENT") == "" {
		panic("TOGGLE_USER_AGENT is required")
	}
	if os.Getenv("TOGGL_WORKSPACE_ID") == "" {
		panic("TOGGLE_WORKSPACE_ID is required")
	}
	if os.Getenv("TOGGL_API_TOKEN") == "" {
		panic("TOGGLE_API_TOKEN is required")
	}

	if os.Getenv("OURA_ACCESS_TOKEN") == "" {
		panic("OURA_ACCESS_TOKEN is required")
	}
}

func main() {
	ctx := context.Background()
	e := echo.New()

	// middleware
	e.Use(middleware.Logger())
	e.Use(middleware.Recover())
	e.Use(session.Middleware(sessions.NewCookieStore([]byte(os.Getenv("SESSION_STORE_KEY")))))
	db, err := dynamodb.NewDB(ctx, os.Getenv("DYNAMODB_URL"), os.Getenv("DYNAMODB_REGION"))
	if err != nil {
		log.Fatal("db cannot build")
	}
	e.Use(database.Middleware(db))
	e.Use(user_session.Middleware(db))

	// template
	tmpl := template.New("template")
	tmpl = tmpl.Funcs(template.FuncMap{
		"toHMS": func(time int) string {
			return fmt.Sprintf("%d:%02d:%02d", time/(60*60), (time/60)%60, time%60)
		},
	})
	tmpl = template.Must(tmpl.ParseGlob("public/views/*.html"))
	renderer := &TemplateRenderer{
		templates: tmpl,
	}
	e.Renderer = renderer

	// router
	e.GET(constant.PATH_AUTH_GOOGLE, google.Index)
	e.GET(constant.PATH_AUTH_GOOGLE_CALLBACK, google.Cb)

	e.GET("/diary", diary)
	e.Logger.Fatal(e.Start(":8080"))
}

func diary(c echo.Context) error {
	db, err := database.Get(c)
	if err != nil {
		return c.String(http.StatusInternalServerError, err.Error())
	}

	now := time.Now()
	total := make(map[string]map[string]int)
	for i := 0; i <= 3; i++ {
		date := now.Add(time.Hour * 24 * -1 * time.Duration(i)).Format("2006-01-02")
		projectTimes := fromToggle(date)
		projectTimes["睡眠"] += fromOura(date)
		projectTimes["移動"] += fromMove(db, date)
		projectTimes["訪問"] += fromVisit(db, date)
		total[date] = projectTimes
	}
	return c.Render(http.StatusOK, "diary.html", map[string]interface{}{
		"day":   now.Format("2006-01-02"),
		"total": total,
	})
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

func fromMove(db *dynamodb.DB, date string) int {

}

func fromVisit(db *dynamodb.DB, date string) int {

}
