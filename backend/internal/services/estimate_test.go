package services

import (
	"errors"
	"testing"
)

func TestEstimate(t *testing.T) {
	for _, points := range []int{1, 3, 100} {
		got, err := Estimate(points)
		if err != nil || got.Points != points || got.Hours != points*4 {
			t.Fatalf("Estimate(%d) = %+v, %v", points, got, err)
		}
	}
	for _, points := range []int{-1, 0, 101} {
		if _, err := Estimate(points); !errors.Is(err, ErrInvalidPoints) {
			t.Errorf("Estimate(%d) error = %v", points, err)
		}
	}
}
