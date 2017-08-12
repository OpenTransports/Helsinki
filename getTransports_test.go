package main

import (
	"fmt"
	"testing"

	"github.com/OpenTransports/lib-go/models"
)

func TestTransports(t *testing.T) {
	stops, err := fetchStops(models.Position{Latitude: 60.192059, Longitude: 24.945831}, 2000.)

	if err != nil {
		fmt.Println(err)
		t.Fail()
	}

	fmt.Println(len(stops))
	for _, s := range stops {
		for _, t := range s.Passages {
			fmt.Println(t)
		}
	}
}