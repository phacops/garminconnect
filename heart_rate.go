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

func (gc *Client) HeartRateByDate(date time.Time) (HeartRate, error) {
	response, err := gc.client.Get(GARMIN_CONNECT_URL + "/modern/proxy/wellness-service/wellness/dailyHeartRate?date=" + date.Format("2006-01-02"))

	if err != nil {
		return HeartRate{}, err
	}

	defer response.Body.Close()

	var heartRate HeartRate

	json.NewDecoder(response.Body).Decode(&heartRate)

	return heartRate, nil
}
