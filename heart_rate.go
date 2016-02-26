package garminconnect

import (
	"encoding/json"
	"time"
)

type HeartRateValue [2]int

type HeartRate struct {
	Min    int              `json:"minHeartRate"`
	Max    int              `json:"maxHeartRate"`
	Values []HeartRateValue `json:"heartRateValues"`
}

func (gc *GarminConnect) HeartRateByDate(date time.Time) HeartRate {
	response, err := gc.client.Get("https://connect.garmin.com/modern/proxy/wellness-service/wellness/dailyHeartRate/" + gc.displayName + "?date=" + date.Format("2006-01-02"))

	if err != nil {
		panic(err)
	}

	defer response.Body.Close()

	var heartRate HeartRate

	json.NewDecoder(response.Body).Decode(&heartRate)

	return heartRate
}
