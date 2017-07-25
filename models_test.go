package main

import (
	"math"
	"testing"
)

// Earth circumference in km
const (
	earthCircumference = 40075
	metersByDegree     = (earthCircumference / 360) * 1000
)

func TestDistanceFrom(t *testing.T) {
	p0 := position{48.82, 2.33}
	p1 := position{48.83, 2.33}
	p2 := position{48.82, 2.34}

	dist := p0.DistanceFrom(&p1)
	if math.Abs(1-dist/1112) > 0.01 {
		t.Fail()
	}

	dist = p0.DistanceFrom(&p2)
	if math.Abs(1-dist/732) > 0.01 {
		t.Fail()
	}

	dist = p0.DistanceFrom(&p0)
	if dist != 0 {
		t.Fail()
	}
}

func TestTransportDistanceFrom(t *testing.T) {
	tr := transport{Position: position{0, 0}}
	p := position{0, 1}
	// Compute distance between the two points
	dist := tr.DistanceFrom(&p)
	// Test that the error is less than 1%
	err := math.Abs(1 - dist/metersByDegree)
	if err > 0.01 {
		t.Fail()
	}
}
