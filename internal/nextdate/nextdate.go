package nextdate

import (
	"fmt"
	"strconv"
	"strings"
	"time"
)

const (
	DateFormat = "20060102"
)

func nextDaily(now, date time.Time, repeat string) (string, error) {
	args := strings.Split(repeat, " ")
	if len(args) != 2 {
		return "", fmt.Errorf("uncorrect repeat format")
	}

	v, err := strconv.Atoi(args[1])
	if err != nil {
		return "", fmt.Errorf("uncorrect repeat format")
	}

	if v < 1 || v > 400 {
		return "", fmt.Errorf("uncorrect repeat format")
	}

	date = date.AddDate(0, 0, v)
	for now.After(date) {
		date = date.AddDate(0, 0, v)
	}

	return date.Format(DateFormat), nil
}

func nextYearly(now, date time.Time) (string, error) {
	date = date.AddDate(1, 0, 0)
	for now.After(date) {
		date = date.AddDate(1, 0, 0)
	}

	return date.Format(DateFormat), nil
}

func NextDate(now time.Time, date string, repeat string) (string, error) {
	if repeat == "" {
		return "", fmt.Errorf("repeat is empty")
	}

	d, err := time.Parse(DateFormat, date)
	if err != nil {
		return "", fmt.Errorf("cannot parse date")
	}

	switch {
	case strings.HasPrefix(repeat, "d"):
		return nextDaily(now, d, repeat)
	case repeat == "y":
		return nextYearly(now, d)
	}

	return "", fmt.Errorf("unexpected type")
}
