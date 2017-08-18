package main

import (
	"fmt"
	"testing"

	"github.com/OpenTransports/lib-go/models"
)

func TestTransports(t *testing.T) {
	transports, err := fetchTransports(
		models.Position{
			Latitude:  60.1665,
			Longitude: 24.9679,
		},
		200.,
	)

	if err != nil {
		fmt.Println(err)
		t.Fail()
	}

	fmt.Println(len(transports))
	for _, trans := range transports {
		fmt.Println(trans)
	}
}
