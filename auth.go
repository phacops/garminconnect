package garminconnect

import (
	"fmt"
	"io/ioutil"
	"net/http"
	"net/http/cookiejar"
	"net/url"
	"regexp"
)

type GarminConnect struct {
	client      *http.Client
	displayName string
}

func NewClient() *GarminConnect {
	cookies, err := cookiejar.New(nil)

	if err != nil {
		panic(err)
	}

	return &GarminConnect{
		client: &http.Client{
			Jar: cookies,
		},
	}
}

func (gc *GarminConnect) Auth(username, password string) bool {
	params := url.Values{}
	params.Set("service", "https://connect.garmin.com/post-auth/login")
	params.Set("clientId", "GarminConnect")
	params.Set("consumeServiceTicket", "false")

	response, err := gc.client.Get("https://sso.garmin.com/sso/login?" + params.Encode())

	if err != nil {
		panic(err)
	}

	defer response.Body.Close()

	if response.StatusCode != http.StatusOK {
		fmt.Println("pre auth request not good")
		return false
	}

	body, err := ioutil.ReadAll(response.Body)

	if err != nil {
		panic(err)
	}

	re := regexp.MustCompile(`name="lt"\s+value="([^\"]+)"`)
	lt := re.FindAllStringSubmatch(string(body), 1)[0][1]

	data := url.Values{}
	data.Set("username", username)
	data.Set("password", password)
	data.Set("_eventId", "submit")
	data.Set("embed", "true")
	data.Set("lt", lt)

	response, err = gc.client.PostForm("https://sso.garmin.com/sso/login?"+params.Encode(), data)

	if err != nil {
		panic(err)
	}

	defer response.Body.Close()

	if response.StatusCode != http.StatusOK {
		fmt.Println("auth request not good")
		return false
	}

	body, err = ioutil.ReadAll(response.Body)

	if err != nil {
		panic(err)
	}

	re = regexp.MustCompile(`ticket=([^']+)'`)
	ticket := re.FindAllStringSubmatch(string(body), 1)[0][1]

	response, err = gc.client.Get("https://connect.garmin.com/post-auth/login?ticket=" + ticket)

	if err != nil {
		panic(err)
	}

	defer response.Body.Close()

	gc.displayName = gc.UserProfile().DisplayName

	return response.StatusCode == http.StatusOK
}
