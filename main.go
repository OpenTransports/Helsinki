package main

import (
	"os"

	"github.com/go-siris/siris"
)

func main() {
	// Create new siris app
	app := siris.New()

	// Serve medias files
	app.StaticWeb("/medias", "./medias")

	// Build api
	// /transports?latitude=...&longitude=...
	app.Get("/transports", GetTransports)
	// /agencies?latitude=...&longitude=...
	app.Get("/agencies", GetAgencies)

	// Set listening port
	var port = os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}

	// Run server
	app.Run(siris.Addr(":"+port), siris.WithCharset("UTF-8"))
}
