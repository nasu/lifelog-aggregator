package toggl

import (
	"encoding/base64"
	"encoding/json"
	"net/http"
	"strconv"
	"time"
)

const (
	details_endpoint    = "https://api.track.toggl.com/reports/api/v2/details"
	password            = "api_token"
	per_page_prediction = 50
)

type Client struct {
	UserAgent   string
	WorkSpaceID string
	ApiToken    string
}

type details struct {
	Id         int          `json:"id"`
	TotalCount int          `json:"total_count"`
	PerPage    int          `json:"per_page"`
	Data       []DetailData `json:"data"`
}

type DetailData struct {
	Id              int       `json:"id"`
	Project         string    `json:"project"`
	Description     string    `json:"description"`
	ProjectHexColor string    `json:"project_hex_color"`
	Start           time.Time `json:"start"`
	End             time.Time `json:"end"`
	Duration        int       `json:"dur"` // unit: milliseconds
}

func (c *Client) GetDetails(since, until string) ([]DetailData, error) {
	page := 1
	data := make([]DetailData, 0, per_page_prediction)
	for {
		d, next, err := c.getDetails(since, until, page)
		if err != nil {
			return nil, err
		}
		data = append(data, d...)
		if !next {
			break
		}
		page++
	}
	return data, nil
}

func (c *Client) getDetails(since, until string, page int) (data []DetailData, next bool, err error) {
	req, err := http.NewRequest("GET", details_endpoint, nil)
	if err != nil {
		return
	}

	q := req.URL.Query()
	q.Add("user_agent", c.UserAgent)
	q.Add("workspace_id", c.WorkSpaceID)
	q.Add("since", since)
	q.Add("until", until)
	q.Add("page", strconv.Itoa(page))
	req.URL.RawQuery = q.Encode()

	req.Header.Add("Authorization", c.authorization())

	body, err := send(req)
	if err != nil {
		return
	}

	var det details
	err = json.Unmarshal(body, &det)
	if err != nil {
		return
	}
	if det.TotalCount > det.PerPage*page {
		next = true
	}
	data = det.Data
	return
}

func (c *Client) authorization() string {
	return "Basic " + base64.StdEncoding.EncodeToString([]byte(c.ApiToken+":"+password))
}
