package geocode

import (
	"context"
	"encoding/json"
	"github.com/androzes/CovidCases/internal/db"
	"io/ioutil"
	"os"
)

type State struct {
	Name      string `json:"name"`
	StateCode string `json:"state_code"`
}

func getStates() ([]State, error) {
	jsonFile, err := os.Open("states.json")
	if err != nil {
		return nil, err
	}

	defer jsonFile.Close()

	byteValue, _ := ioutil.ReadAll(jsonFile)
	var states []State
	json.Unmarshal([]byte(byteValue), &states)
	return states, nil
}

func UpdateStates(ctx context.Context) error {
	states, err := getStates()
	if err != nil {
		return err
	}

	covidStates := []db.CovidStat{}
	for _, state := range states {
		covidStates = append(covidStates, db.CovidStat{
			Code: state.StateCode,
			Name: state.Name,
		})
	}

	return db.UpdateStates(ctx, covidStates)
}