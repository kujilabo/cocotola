package gateway

import "time"

var jst *time.Location

func init() {
	tz, err := time.LoadLocation("Asia/Tokyo")
	if err != nil {
		panic(err)
	}
	jst = tz
}
