package main

import (
	"os"

	"github.com/OpenTransports/lib-go/models"
	"github.com/go-siris/siris/context"
)

var serverURL = os.Getenv("SERVER_URL")

// HSL agency description
var HSL = models.Agency{
	ID:   "Finland.Helsinki.HSL",
	Name: "HSL",
	URL:  "https://hsl.fi",
	Center: models.Position{
		Latitude:  60.192059,
		Longitude: 24.945831,
	},
	Radius: 20000,
	Types: map[int]models.TypeInfo{
		models.Tram: models.TypeInfo{
			Name: models.TramString,
			Icon: serverURL + "/medias/tram.png",
		},
		models.Metro: models.TypeInfo{
			Name: models.MetroString,
			Icon: serverURL + "/medias/metro.png",
		},
		models.Rail: models.TypeInfo{
			Name: models.RailString,
			Icon: serverURL + "/medias/train.png",
		},
		models.Bus: models.TypeInfo{
			Name: models.BusString,
			Icon: serverURL + "/medias/bus.png",
		},
		models.Ferry: models.TypeInfo{
			Name: models.FerryString,
			Icon: serverURL + "/medias/ferry.png",
		},
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
