package garminconnect

import (
	"encoding/json"
	"time"
)

type HeartRateValue [2]int

type DailyHeartRate struct {
	Min    int              `json:"minHeartRate"`
	Max    int              `json:"maxHeartRate"`
	Values []HeartRateValue `json:"heartRateValues"`
}

func (gc *GarminConnect) DailyHeartRate(date time.Time) DailyHeartRate {
	response, err := gc.client.Get("https://connect.garmin.com/modern/proxy/wellness-service/wellness/dailyHeartRate/" + gc.displayName + "?date=" + date.Format("2006-01-02"))

	if err != nil {
		panic(err)
	}

	defer response.Body.Close()

	var dailyHeartRate DailyHeartRate

	json.NewDecoder(response.Body).Decode(&dailyHeartRate)

	return dailyHeartRate
}
