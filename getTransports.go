package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"strconv"
	"time"

	"github.com/OpenTransports/lib-go/models"
	"github.com/go-siris/siris/context"
)

// GetTransports - /transports?latitude=...&longitude=...&radius=...
// Send the transports aroud the passed position
// @formParam latitude : optional, the latitude around where to search, default is 0
// @formParam longitude : optional, the longitude around where to search, default is 0
// @formParam radius : optional, default is 200m
func GetTransports(ctx context.Context) {
	ctx.Header("Access-Control-Allow-Origin", "*")
	// Get position in params
	// Parse them to floats
	// Ignore errors because it default to 0
	latitude, _ := strconv.ParseFloat(ctx.FormValue("latitude"), 64)
	longitude, _ := strconv.ParseFloat(ctx.FormValue("longitude"), 64)
	radius, errRadius := strconv.ParseFloat(ctx.FormValue("radius"), 64)
	// Create a Position object
	position := models.Position{
		Latitude:  latitude,
		Longitude: longitude,
	}
	// Set the radius to its default value if none is passed
	if errRadius != nil {
		radius = 200.0
	}
	// Make query
	stops, err := fetchStops(position, radius)
	if err != nil {
		ctx.Application().Log("Error making query answer in /transports\n	==> %v", err)
	}
	// Write results
	_, err = ctx.JSON(stops)
	// Log the error if any
	if err != nil {
		ctx.Application().Log("Error writting answer in /transports\n	==> %v", err)
	}
}

// Fetch the nearby transports from the HSL server
// https://digitransit.fi/en
// http://dev.hsl.fi/graphql/console
func fetchStops(userPosition models.Position, radius float64) ([]models.Transport, error) {
	// Build query and make request
	query := fmt.Sprintf(queryTemplate, userPosition.Latitude, userPosition.Longitude, radius)
	const URL = "http://api.digitransit.fi/routing/v1/routers/finland/index/graphql"
	req, err := http.NewRequest("POST", URL, bytes.NewBuffer([]byte(query)))
	if err != nil {
		return nil, err
	}
	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()
	// Parse answer and put it in a struct
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}
	answer := &answerType{}
	err = json.Unmarshal(body, answer)
	if err != nil {
		return nil, err
	}
	// Map anwer to a curated array of nearby transports
	// The answer is not currectly structured
	// A stop (node) containes all the futur passage for each line
	// We need to flatten the answer, then merge what can be merged and remove what is useless
	// 3. Filter out furthest transports for each lines
	fmt.Println(len(reduceStops(
		// 1. Flatten the answer and map to transports structs
		mapStops(answer),
	)))
	return filterStops(
		// 2. Merge stops with the same Name and Line
		reduceStops(
			// 1. Flatten the answer and map to transports structs
			mapStops(answer),
		),
		userPosition,
	), nil
}

// 1. Flatten the answer and map to transports structs
func mapStops(answer *answerType) []models.Transport {
	stops := []models.Transport{}
	for _, e := range answer.Data.Nearest.Edges {
		t := e.Node.Place
		for _, time := range t.Times {
			stops = append(stops, models.Transport{
				ID:   t.ID,
				Name: t.Name,
				Line: time.Trip.Route.Name,
				Type: modeToType(time.Trip.Route.Mode),
				Position: models.Position{
					Latitude:  t.Latitude,
					Longitude: t.Longitude,
				},
				Passages: []models.Passage{
					models.Passage{
						Direction: time.Trip.Direction,
						Times:     []string{absoluteDateToRelativeDate(time.Date)},
					},
				},
			})
		}
	}
	return stops
}

// 2. Merge stops with the same Name and Line
func reduceStops(stops []models.Transport) []models.Transport {
	reducedStops := []models.Transport{}
	for _, stop1 := range stops {
		added := false
		for _, stop2 := range reducedStops {
			if stop1.Name == stop2.Name && stop1.Line == stop2.Line {
				added = true
				stop2.Passages = mergePassages(stop1.Passages, stop2.Passages)
			}
		}
		if !added {
			reducedStops = append(reducedStops, stop1)
		}
	}
	return reducedStops
}

func mergePassages(passages1 []models.Passage, passages2 []models.Passage) []models.Passage {
	for _, p1 := range passages1 {
		added := false
		for _, p2 := range passages2 {
			if p1.Direction == p2.Direction {
				p2.Times = append(p2.Times, p1.Times...)
				added = true
				break
			}
		}
		if !added {
			passages2 = append(passages2, p1)
		}
	}

	return passages2
}

// 3. Filter out furthest transports for each lines
func filterStops(stops []models.Transport, userPosition models.Position) []models.Transport {
	filterdStops := []models.Transport{}
	for _, stop1 := range stops {
		added := false
		for i, stop2 := range filterdStops {
			if stop1.Line == stop2.Line {
				if stop1.Position.DistanceFrom(userPosition) < stop2.Position.DistanceFrom(userPosition) {
					filterdStops[i] = stop1
				}
				added = true
			}
		}
		if !added {
			filterdStops = append(filterdStops, stop1)
		}
	}
	return filterdStops
}

// From a date (int) representing the number of seconds since midnight,
// return the number of minute between now and the date
func absoluteDateToRelativeDate(date int) string {
	now := time.Now()
	nowSec := (now.Hour()*60 + now.Minute()) * 60
	return fmt.Sprintf("%v mn", (date-nowSec)/60)
}

// Map the stringifyed transports type to our type
func modeToType(mode string) int {
	switch mode {
	case "TRAM":
		return models.Tram
	default:
		return models.Unknown
	}
}
