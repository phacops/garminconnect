package garminconnect

import (
	"encoding/json"
	"net/url"
	"time"
)

type Sleep struct {
	Date       string `json:"calendarDate"`
	Duration   int64  `json:"sleepTimeSeconds"`
	BedTime    int64  `json:"sleepStartTimestampGMT"`
	WakeUpTime int64  `json:"sleepEndTimestampGMT"`
}

func (gc *Client) SleepByDate(date time.Time) (Sleep, error) {
	params := url.Values{}
	params.Set("date", date.Format("2006-01-02"))
	params.Set("nonSleepBufferMinutes", "60")

	response, err := gc.client.Get(GARMIN_CONNECT_URL + "/modern/proxy/wellness-service/wellness/dailySleep?" + params.Encode())

	if err != nil {
		return Sleep{}, err
	}

	defer response.Body.Close()

	var sleep Sleep

	json.NewDecoder(response.Body).Decode(&sleep)

	return sleep, nil
}
