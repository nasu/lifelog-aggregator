package toggl

import (
	"encoding/json"
	"errors"
	"io/ioutil"
	"log"
	"net/http"
)

type errorObject struct {
	Error struct {
		Message string
	}
}

func send(req *http.Request) ([]byte, error) {
	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		log.Printf("toggl request: status=%d", resp.StatusCode)
		b, err := ioutil.ReadAll(resp.Body)
		if err != nil {
			return nil, err
		}
		var eo errorObject
		err = json.Unmarshal(b, &eo)
		if err != nil {
			return nil, err
		}
		return nil, errors.New("toggl request: " + eo.Error.Message)
	}
	return ioutil.ReadAll(resp.Body)
}
