package oura

import (
	"encoding/json"
	"time"
)

const (
	userinfo_endpoint      = "https://api.ouraring.com/v1/userinfo"
	sleep_endpoint         = "https://api.ouraring.com/v1/sleep"
	activity_endpoint      = "https://api.ouraring.com/v1/activity"
	readiness_endpoint     = "https://api.ouraring.com/v1/readiness"
	ideal_bedtime_endpoint = "https://api.ouraring.com/v1/bedtime"
)

type Client struct {
	AccessToken string
}

type UserInfo struct {
	Age    int     `json:"age"`
	Weight float32 `json:"weight"`
	Gender string  `json:"gender"`
	Email  string  `json:"email"`
}

type Sleep struct {
	Sleep []*SleepOneDay `json:"sleep"`
}

type SleepOneDay struct {
	SummaryDate       string    `json:"summary_date"`
	PeriodID          int       `json:"period_id"`
	IsLongest         int       `json:"is_longest"`
	Timezone          int       `json:"timezone"`
	BedtimeStart      time.Time `json:"bedtime_start"`
	BedtimeEnd        time.Time `json:"bedtime_end"`
	Score             int       `json:"score"`
	ScoreTotal        int       `json:"score_total"`
	ScoreDisturbances int       `json:"score_disturbances"`
	ScoreEfficiency   int       `json:"score_efficiency"`
	ScoreLatency      int       `json:"score_latency"`
	ScoreRem          int       `json:"score_rem"`
	ScoreDeep         int       `json:"score_deep"`
	ScoreAlignment    int       `json:"score_alignment"`
	Total             int       `json:"total"`    // duration - awake
	Duration          int       `json:"duration"` // from bed-in to bed-out
	Awake             int       `json:"awake"`
	Light             int       `json:"light"`
	Rem               int       `json:"rem"`
	Deep              int       `json:"deep"`
	OnsetLatency      int       `json:"onset_latency"`
	Restless          int       `json:"restless"`
	Efficiency        int       `json:"efficiency"`
	MidpointTime      int       `json:"midpoint_time"`
	HRLowest          float32   `json:"hr_lowest"`
	HRAverage         float32   `json:"hr_average"`
	Rmssd             int       `json:"rmssd"`
	BreathAverage     float32   `json:"breath_average"`
	TemperatureDelta  float32   `json:"temperature_delta"`
	Hypnogram5min     string    `json:"hypnogram_5min"`
	HR5min            []int     `json:"hr_5min"`
	Rmssd5min         []int     `json:"rmssd_5min"`
}

type Readiness struct {
	Readiness []*ReadinessOneDay `json:"readiness"`
}

type ReadinessOneDay struct {
	SummaryDate          string `json:"summary_date"`
	PeriodID             int    `json:"period_id"`
	Score                int    `json:"score"`
	ScorePreviousNight   int    `json:"score_previous_night"`
	ScoreSleepBalance    int    `json:"score_sleep_balance"`
	ScorePreviousDay     int    `json:"score_previous_day"`
	ScoreActivityBalance int    `json:"score_activity_balance"`
	ScoreRestingHR       int    `json:"socre_resting_hr"`
	ScoreHRVBalance      int    `json:"score_hrv_balance"`
	ScoreRecoveryIndex   int    `json:"score_recovery_index"`
	ScoreTemperature     int    `json:"score_temperature"`
	RestModeState        int    `json:"score_rest_mode_state"`
}

func (c *Client) UserInfo() (res *UserInfo, err error) {
	params := map[string]string{}
	body, err := send("GET", userinfo_endpoint, params, c.authorization())
	if err != nil {
		return
	}

	err = json.Unmarshal(body, &res)
	if err != nil {
		return
	}
	return
}

func (c *Client) ReadinessOneDay(date string) (res *ReadinessOneDay, err error) {
	params := map[string]string{"start": date, "end": date}
	body, err := send("GET", readiness_endpoint, params, c.authorization())
	if err != nil {
		return
	}

	var readiness Readiness
	err = json.Unmarshal(body, &readiness)
	if err != nil {
		return
	}
	if len(readiness.Readiness) == 0 {
		res = &ReadinessOneDay{}
	} else {
		res = readiness.Readiness[0]
	}
	return
}

func (c *Client) SleepOneDay(date string) (res *SleepOneDay, err error) {
	params := map[string]string{"start": date, "end": date}
	body, err := send("GET", sleep_endpoint, params, c.authorization())
	if err != nil {
		return
	}

	var sleep Sleep
	err = json.Unmarshal(body, &sleep)
	if err != nil {
		return
	}
	if len(sleep.Sleep) == 0 {
		res = &SleepOneDay{}
	} else {
		res = sleep.Sleep[0]
	}
	return
}

func (c *Client) authorization() string {
	return "Bearer " + c.AccessToken
}
