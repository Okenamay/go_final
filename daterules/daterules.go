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

func NextTime(now time.Time, date string, repeat string) (string, error) {
	if repeat == "" {
		return "", errors.New("no repeat parameter")
	}

	startDate, err := time.Parse(TimeFormat, date)
	if err != nil {
		return "", err
	}

	parts := strings.Fields(repeat)

	if !strings.Contains("yd", parts[0]) {
		return "", errors.New("wrong repeat format")
	} else if parts[0] == "d" && len(parts) != 2 {
		return "", errors.New("wrong repeat time format")
	} else if parts[0] == "y" && len(parts) != 1 {
		return "", errors.New("wrong daily repeat format")
	}

	for {
		if parts[0] == "y" {
			startDate = startDate.AddDate(1, 0, 0)
		} else if parts[0] == "d" {
			part, err := strconv.Atoi(parts[1])
			if part > 366 || part <= 0 || err != nil {
				return "", errors.New("wrong repeat time")
			}
			startDate = startDate.AddDate(0, 0, part)
		}

		if startDate.After(now) || startDate.Equal(now) {
			return startDate.Format(TimeFormat), nil
		}
	}
}
