package services

import "errors"

var ErrInvalidPoints = errors.New("points must be an integer between 1 and 100")

type Estimation struct {
	Points int `json:"points"`
	Hours  int `json:"hours"`
}

func Estimate(points int) (Estimation, error) {
	if points < 1 || points > 100 {
		return Estimation{}, ErrInvalidPoints
	}
	return Estimation{Points: points, Hours: points * 4}, nil
}
