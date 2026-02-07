package utils

import (
	"os"
	"time"
)

var AppLocation *time.Location

func InitTimezone() {
	tz := os.Getenv("APP_TIMEZONE")
	if tz == "" {
		tz = "Asia/Jakarta"
	}

	loc, err := time.LoadLocation(tz)
	if err != nil {
		panic(err)
	}

	AppLocation = loc
}
