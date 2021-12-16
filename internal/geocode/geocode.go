package geocode

import (
	"fmt"
	"encoding/json"
	"io/ioutil"
	"net/http"
	"os"
)

type Response struct {
	Address Address `json:"address"`
}

type Address struct {
	State string `json:"state"`
	Country string `json:"country"`
}

const API_HOST = "https://us1.locationiq.com/v1/reverse.php"

func GetPlaceDetails(latitude string, longitude string) (Response, error) {
	client := &http.Client{}

	req, err := http.NewRequest("GET", API_HOST, nil)
	if err != nil {
		return Response{}, err
	}

	apiKey := os.Getenv("LOCATIONIQ_API_KEY")

	q := req.URL.Query()
	q.Add("key", apiKey)
	q.Add("lat", latitude)
	q.Add("lon", longitude)
	q.Add("format", "json")
	req.URL.RawQuery = q.Encode()

	resp, err := client.Do(req)
	if err != nil {
		return Response{}, err
	}

	defer resp.Body.Close()
	bodyBytes, _ := ioutil.ReadAll(resp.Body)

	fmt.Println("Response: "+ string(bodyBytes))

	// Convert response body to Response
	var responseObj Response
	err = json.Unmarshal(bodyBytes, &responseObj)
	if err != nil {
		return Response{}, err
	}

	return responseObj, nil
}
