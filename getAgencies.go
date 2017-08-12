package main

import (
	"github.com/OpenTransports/lib-go/models"
	"github.com/go-siris/siris/context"
)

// HSL agency description
var HSL = models.Agency{
	ID:     "Finland.Helsinki.HSL",
	Name:   "HSL",
	URL:    "https://hsl.fi",
	Git:    "https://github.com/OpenTransports/Helsinki",
	Center: models.Position{Latitude: 60.192059, Longitude: 24.945831},
	Radius: 20000,
	Types: []int{
		models.Tram,
		models.Metro,
		models.Rail,
		models.Bus,
		models.Ferry,
	},
	TypesString: []string{
		models.TramString,
		models.MetroString,
		models.RailString,
		models.BusString,
		models.FerryString,
	},
}

// GetAgencies - /api/agencies
// Send the agencies handled by this server
func GetAgencies(ctx context.Context) {
	ctx.Header("Access-Control-Allow-Origin", "*")
	_, err := ctx.JSON([]models.Agency{HSL})
	// Log the error if any
	if err != nil {
		ctx.Application().Log("Error writting answer in /api/agencies\n	==> %v", err)
	}
}
