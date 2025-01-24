package daterules

import (
	"errors"
	"strconv"
	"strings"
	"time"
)

const TimeFormat string = "20060102"

type Task struct {
	ID      string `json:"id"`
	Date    string `json:"date"`
	Title   string `json:"title"`
	Comment string `json:"comment"`
	Repeat  string `json:"repeat"`
}

func NextTime(now time.Time, date string, repeat string) (time.Time, error) {
	if repeat == "" {
		return now, nil
	}

	parts := strings.Fields(repeat)

	if !strings.Contains("yd", parts[0]) {
		return time.Time{}, errors.New("wrong repeat format")
	} else if parts[0] == "d" && len(parts) != 2 {
		return time.Time{}, errors.New("wrong repeat time format")
	} else if parts[0] == "y" && len(parts) != 1 {
		return time.Time{}, errors.New("wrong yearly repeat format")
	}

	for {
		if parts[0] == "y" {
			return now.AddDate(1, 0, 0), nil
		} else if parts[0] == "d" {
			part, err := strconv.Atoi(parts[1])
			if part > 366 || part <= 0 || err != nil {
				return time.Time{}, errors.New("wrong repeat time")
			}
			return now.AddDate(0, 0, part), nil
		}
	}

}

func (t *Task) NextDay() (string, error) {
	currentTime, err := time.Parse(TimeFormat, t.Date)
	if err != nil {
		return "", err
	}

	nT, err := NextTime(currentTime, "", t.Repeat)
	if err != nil {
		return "", err
	}

	return nT.Format(TimeFormat), nil
}
