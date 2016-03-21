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

func (gc *Client) SleepByDate(date time.Time) Sleep {
	params := url.Values{}
	params.Set("date", date.Format("2006-01-02"))
	params.Set("nonSleepBufferMinutes", "60")

	response, err := gc.client.Get("https://connect.garmin.com/modern/proxy/wellness-service/wellness/dailySleep/user/" + gc.displayName + "?" + params.Encode())

	if err != nil {
		panic(err)
	}

	defer response.Body.Close()

	var sleep Sleep

	json.NewDecoder(response.Body).Decode(&sleep)

	return sleep
}
