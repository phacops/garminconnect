package garminconnect

import "encoding/json"

type UserProfile struct {
	DisplayName string `json:"displayName"`
}

func (gc *Client) UserProfile() UserProfile {
	response, err := gc.client.Get("http://connect.garmin.com/proxy/userprofile-service/socialProfile")

	if err != nil {
		panic(err)
	}

	defer response.Body.Close()

	var profile UserProfile

	json.NewDecoder(response.Body).Decode(&profile)

	return profile
}
