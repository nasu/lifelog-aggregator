package main

import (
	"fmt"
	"log"
	"os"
	"time"

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
	now := time.Now()
	total := make(map[string]map[string]int)
	for i := 1; i <= 7; i++ {
		date := now.Add(time.Hour * 24 * -1 * time.Duration(i)).Format("2006-01-02")
		projectTimes := fromToggle(date)
		projectTimes["睡眠"] += fromOura(date)
		total[date] = projectTimes
	}
	fmt.Println(total)
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
