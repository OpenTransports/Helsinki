package main

import (
	"github.com/go-siris/siris/context"
)

// HSL agency description
var HSL = agency{
	ID:     "Finland.Helsinki.HSL",
	Name:   "HSL",
	URL:    "https://hsl.fi",
	Git:    "https://github.com/OpenTransports/Helsinki",
	Center: position{Latitude: 60.192059, Longitude: 24.945831},
	Radius: 20000,
	Types: []int{
		tram,
		metro,
		rail,
		bus,
		ferry,
	},
	TypesString: []string{
		tramString,
		metroString,
		railString,
		busString,
		ferryString,
	},
}

// GetAgencies - /api/agencies
// Send the agencies handled by this server
func GetAgencies(ctx context.Context) {
	ctx.Header("Access-Control-Allow-Origin", "*")
	_, err := ctx.JSON([]agency{HSL})
	// Log the error if any
	if err != nil {
		ctx.Application().Log("Error writting answer in /api/agencies\n	==> %v", err)
	}
}
