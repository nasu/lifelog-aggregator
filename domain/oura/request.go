package oura

import (
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
)

type errorObject struct {
	Status int
	Title  string
	Detail string
}

func send(method, endpoint string, params map[string]string, auth string) (body []byte, err error) {
	req, err := http.NewRequest(method, endpoint, nil)
	req.Header.Add("Authorization", auth)

	if err != nil {
		return
	}
	q := req.URL.Query()
	for k, v := range params {
		q.Add(k, v)
	}
	req.URL.RawQuery = q.Encode()

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		log.Printf("oura request: status=%d", resp.StatusCode)
		b, err := ioutil.ReadAll(resp.Body)
		if err != nil {
			return nil, err
		}
		var eo errorObject
		err = json.Unmarshal(b, &eo)
		if err != nil {
			return nil, err
		}
		return nil, errors.New(fmt.Sprintf("oura request: title=%s detail=%s", eo.Title, eo.Detail))
	}
	return ioutil.ReadAll(resp.Body)
}
