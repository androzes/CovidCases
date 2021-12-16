package main

import (
	"context"
	"fmt"
	"github.com/androzes/CovidCases/internal/covid"
	"github.com/androzes/CovidCases/internal/db"
	"github.com/androzes/CovidCases/internal/geocode"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	"net/http"
	"strings"
)

func main() {
	// Echo instance
	e := echo.New()

	// Middleware
	e.Use(middleware.Logger())
	e.Use(middleware.Recover())

	// Routes
	e.GET("/covid/user/:lat_lon", getUserLocationCovidCases)

	e.GET("/user", getUserPlace)

	e.GET("/covid/state/:code", getCovidDataByStateCode)

	e.POST("/covid/update", updateCovidData)

	// Start server
	e.Logger.Fatal(e.Start(":1323"))
}

// Handler
func getUserLocationCovidCases(c echo.Context) error {
	latLonStr := c.Param("lat_lon")
	coords := strings.Split(latLonStr, ",")

	if len(coords) != 2 {
		return echo.NewHTTPError(http.StatusBadRequest, "Invalid coordinates")
	}

	lat := strings.Trim(coords[0], " ")
	lng := strings.Trim(coords[1], " ")

	if !isLatLongValid(lat, lng) {
		return echo.NewHTTPError(http.StatusBadRequest, "Invalid coordinates")
	}

	// get state and country from coordinates
	response, err := geocode.GetPlaceDetails(lat, lng)
	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, err.Error())
	}

	// check if lat long belong to india or not
	if !strings.EqualFold(response.Address.Country, "india") {
		return echo.NewHTTPError(http.StatusBadRequest, "Coordinates should be within India")
	}

	// check if state is not empty
	if response.Address.State == "" {
		return echo.NewHTTPError(http.StatusBadRequest, "Unable to figure out the state for given coordinates")
	}

	ctx := context.TODO()

	// get covid cases for this country and state
	stateStat, err := db.GetStatsByStateName(ctx, response.Address.State)
	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, err.Error())
	}

	countryStat, err := db.GetTotalStatsForCountry(ctx)
	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, err.Error())
	}

	var result struct {
		State db.CovidStat `json:"state"`
		Country db.CovidStat `json:"country"`
	}

	result.State = stateStat
	result.Country = countryStat

	return c.JSON(http.StatusOK, result)
}

func isLatLongValid(lat string, lng string) bool {
	return lat != "" && lng != "";
}

func getUserPlace(c echo.Context) error {
	lat := c.QueryParam("lat")
	lng := c.QueryParam("lng")

	if !isLatLongValid(lat, lng) {
		return echo.NewHTTPError(http.StatusBadRequest, "Invalid coordinates")
	}

	response, err := geocode.GetPlaceDetails(lat, lng)
	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, err.Error())
	}

	// check if lat long belong to india or not
	if !strings.EqualFold(response.Address.Country, "india") {
		return echo.NewHTTPError(http.StatusBadRequest, "Coordinates should be within India")
	}

	// check if state is not empty
	if response.Address.State == "" {
		return echo.NewHTTPError(http.StatusBadRequest, "Unable to figure out the state for the coordinates")
	}

	return c.JSON(http.StatusOK, response)
}

func updateCovidData(c echo.Context) error {
	ctx := context.TODO()

	err := geocode.UpdateStates(ctx)
	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, err.Error())
	}

	response, err := covid.GetData()
	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, err.Error())
	}

	stats := []db.CovidStat{}

	for stateCode, stateData := range response.Data {
		stats = append(stats, db.CovidStat{
			Code: stateCode,
			NumCovidCases: stateData.Total.Confirmed - stateData.Total.Recovered - stateData.Total.Deceased,
			LastUpdated: stateData.Meta.LastUpdated,
		})
	}

	err = db.UpdateCovidStats(ctx, stats)
	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, err.Error())
	}

	return c.JSON(http.StatusOK, "Done")
}

func getCovidDataByStateCode(c echo.Context) error {
	stateCodeStr := c.Param("code")
	stateCodes := strings.Split(stateCodeStr, ",")

	fmt.Println("state code param: " + stateCodeStr)
	// validate codes
	for _, stateCode := range stateCodes {
		if len(stateCode) <= 0 ||  len(stateCode) > 2 {
			return echo.NewHTTPError(http.StatusBadRequest, "Invalid state code: " + stateCode)
		}
	}

	covidStats, err := db.GetStatsByStateCodes(context.TODO(), stateCodes)
	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, err.Error())
	}

	return c.JSON(http.StatusOK, covidStats)
}
