package covid

import (
	"encoding/json"
	"errors"
	"io/ioutil"
	"net/http"
	"time"
)

type Covid interface {
	GetData() Response
}

type Response struct {
	Data map[string]StateData
}

type StateData struct {
	Meta  Meta  `json:"meta"`
	Total Total `json:"total"`
}

type Meta struct {
	LastUpdated time.Time `json:"last_updated"`
}

type Total struct {
	Confirmed int `json:"confirmed"`
	Deceased  int `json:"deceased"`
	Recovered int `json:"recovered"`
}

func GetData() (Response, error) {
	resp, err := http.Get("https://data.covid19india.org/v4/min/data.min.json")
	if err != nil {
		return Response{}, errors.New(err.Error())
	}

	defer resp.Body.Close()
	bodyBytes, _ := ioutil.ReadAll(resp.Body)

	// Convert response body to Response
	var stateMap map[string]StateData
	err = json.Unmarshal(bodyBytes, &stateMap)
	if err != nil {
		return Response{}, errors.New(err.Error())
	}

	return Response{Data: stateMap}, nil
}
