package util

import (
	"context"
	"errors"
	"strconv"
	"time"
)

func ToDateTime(dtStr string) (time.Time, error) {
	layout := "2006-01-02 15:04:05"
	if dtStr != "" {
		parsedDate, err := time.Parse(layout, dtStr)
		if err != nil {
			return time.Time{}, err
		}
		return parsedDate, nil
	}
	return time.Time{}, nil
}

// Set Difference: A - B
func Difference(a, b []int64) (diff []int64) {
	m := make(map[int64]bool)
	for _, item := range b {
		m[item] = true
	}
	for _, item := range a {
		if _, ok := m[item]; !ok {
			diff = append(diff, item)
		}
	}
	return
}

func GetUserID(ctx context.Context) (int, error) {
	userId := ctx.Value("userId")
	if userId == nil {
		return 0, errors.New("invalid user id")
	} else {
		switch userId.(type) {
		case int:
			return ctx.Value("userId").(int), nil
		case string:
			return strconv.Atoi(ctx.Value("userId").(string))
		default:
			return 0, errors.New("invalid user id")
		}
	}
}
