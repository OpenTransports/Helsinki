package main

import (
	"context"
	"fmt"
	"strconv"
	"time"

	"github.com/OpenTransports/lib-go/models"
	irisContext "github.com/go-siris/siris/context"
	"github.com/shurcooL/graphql"
)

// GetTransports - /transports?latitude=...&longitude=...&radius=...
// Send the transports aroud the passed position
// @formParam latitude : optional, the latitude around where to search, default is 0
// @formParam longitude : optional, the longitude around where to search, default is 0
// @formParam radius : optional, default is 200m
func GetTransports(ctx irisContext.Context) {
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
	transports, err := fetchTransports(position, radius)
	if err != nil {
		ctx.Application().Log("Error making query answer in /transports\n	==> %v", err)
	}
	// Write results
	_, err = ctx.JSON(transports)
	// Log the error if any
	if err != nil {
		ctx.Application().Log("Error writting answer in /transports\n	==> %v", err)
	}
}

// Fetch the nearby transports from the HSL server
// https://digitransit.fi/en
// http://dev.hsl.fi/graphql/console
var client = graphql.NewClient("http://api.digitransit.fi/routing/v1/routers/hsl/index/graphql", nil, nil)

func fetchTransports(userPosition models.Position, radius float64) ([]models.Transport, error) {
	// Make query
	answer := &queryStruct{}
	variables := map[string]interface{}{
		"lat":    graphql.Float(userPosition.Latitude),
		"lon":    graphql.Float(userPosition.Longitude),
		"radius": graphql.Int(radius),
	}
	err := client.Query(context.Background(), answer, variables)
	if err != nil {
		fmt.Println(err)
	}
	// Map anwer to a curated array of nearby transports
	// The answer is not currectly structured
	// A transport (node) containes all the futur passage for each line
	// We need to flatten the answer, then merge what can be merged and remove what is useless
	// 3. Filter out furthest transports for each lines
	return filterDistantTransports(
		// 2. Merge transports with the same Name and Line
		reduceTransports(
			// 1. Flatten the answer and map to transports structs
			mapToTransports(answer),
		),
		userPosition,
	), nil
}

// 1. Flatten the answer and map to transports structs
func mapToTransports(answer *queryStruct) []models.Transport {
	transports := []models.Transport{}
	for _, e := range answer.Nearest.Edges {
		stop := e.Node.Stop
		for _, passage := range stop.StoptimesWithoutPatterns {
			transports = append(transports, models.Transport{
				ID:       stop.GtfsID,
				AgencyID: "Finland.Helsinki.HSL",
				Name:     passage.Trip.Route.ShortName + " - " + stop.Name,
				Line:     passage.Trip.Route.ShortName,
				Type:     modeToType(passage.Trip.Route.Mode),
				Position: models.Position{
					Latitude:  stop.Lat,
					Longitude: stop.Lon,
				},
				Informations: []models.Information{
					models.Information{
						Title:     passage.Trip.TripHeadsign,
						Content:   []string{absoluteDateToRelativeDate(passage.RealtimeDeparture)},
						Timestamp: int(time.Now().Unix() * 1000),
					},
				},
			})
		}
	}
	return transports
}

// 2. Merge transport with the same Name and Line
func reduceTransports(transports []models.Transport) []models.Transport {
	reducedTransports := []models.Transport{}
	for _, transport1 := range transports {
		added := false
		for i, transport2 := range reducedTransports {
			if transport1.Name == transport2.Name && transport1.Line == transport2.Line {
				added = true
				reducedTransports[i].Informations = mergeInformations(transport2.Informations, transport1.Informations)
			}
		}
		if !added {
			reducedTransports = append(reducedTransports, transport1)
		}
	}
	return reducedTransports
}

// Merge Informations into informations1
func mergeInformations(informations1 []models.Information, informations2 []models.Information) []models.Information {
	for _, info2 := range informations2 {
		added := false
		for i, info1 := range informations1 {
			if info1.Title == info2.Title {
				informations1[i].Content = append(info1.Content, info2.Content...)
				added = true
				break
			}
		}
		if !added {
			informations1 = append(informations1, info2)
		}
	}

	return informations1
}

// 3. Filter out furthest transports for each lines
func filterDistantTransports(transports []models.Transport, userPosition models.Position) []models.Transport {
	filterdTransports := []models.Transport{}
	for _, transport1 := range transports {
		added := false
		for i, transport2 := range filterdTransports {
			if transport1.Line == transport2.Line {
				if transport1.Position.DistanceFrom(userPosition) < transport2.Position.DistanceFrom(userPosition) {
					filterdTransports[i] = transport1
				}
				added = true
			}
		}
		if !added {
			filterdTransports = append(filterdTransports, transport1)
		}
	}
	return filterdTransports
}

// From a date (int) representing the number of seconds since midnight,
// return the number of minute between now and the date
func absoluteDateToRelativeDate(date int) string {
	_, localOffset := time.Now().Zone()
	helsinkiOffset := 60 * 60 * 3
	now := time.Now()
	nowSec := (now.Hour()*60+now.Minute())*60 + (helsinkiOffset - localOffset)
	waitingTime := (date - nowSec) / 60
	if waitingTime < 59 {
		return fmt.Sprintf("%v mn", waitingTime)
	}

	hours := waitingTime / 60
	minutes := waitingTime % 60
	return fmt.Sprintf("%v h %v mn", hours, minutes)
}

func modeToType(mode string) int {
	switch mode {
	case "TRAM":
		return models.Tram
	case "BUS":
		return models.Bus
	case "SUBWAY":
		return models.Metro
	case "RAIL":
		return models.Rail
	case "FERRY":
		return models.Ferry
	default:
		return models.Unknown
	}
}
