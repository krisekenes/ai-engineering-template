package services

import "testing"

func TestIndependentEstimateCap(t *testing.T) {
	for _, points := range []int{1, 3, 20} {
		got, err := Estimate(points)
		if err != nil || got.Hours != points*4 || got.Points != points {
			t.Fatalf("valid points %d: got %+v, error %v", points, got, err)
		}
	}
	for _, points := range []int{-1, 0, 21, 100, 101} {
		if _, err := Estimate(points); err == nil {
			t.Errorf("points %d must be rejected", points)
		}
	}
}
