package garminconnect

import (
	"errors"
	"io/ioutil"
	"net/http"
	"net/http/cookiejar"
	"net/url"
	"regexp"
)

const (
	GARMIN_CONNECT_URL = "https://connect.garmin.com"
)

type Client struct {
	client      *http.Client
	displayName string
}

func NewClient(httpClient ...*http.Client) (*Client, error) {
	cookies, err := cookiejar.New(nil)

	if err != nil {
		return nil, err
	}

	var client *http.Client

	if len(httpClient) > 0 {
		client = httpClient[0]
	} else {
		client = &http.Client{}
	}

	client.Jar = cookies

	return &Client{
		client: client,
	}, nil
}

func (gc *Client) Auth(username, password string) error {
	loginUrl, err := url.Parse("https://sso.garmin.com/sso/login?service=" + GARMIN_CONNECT_URL + "/post-auth/login&webhost=olaxpw-connect13.garmin.com&source=" + GARMIN_CONNECT_URL + "/en-US/signin&redirectAfterAccountLoginUrl=" + GARMIN_CONNECT_URL + "/post-auth/login&redirectAfterAccountCreationUrl=" + GARMIN_CONNECT_URL + "/post-auth/login&gauthHost=https://sso.garmin.com/sso&locale=en_US&id=gauth-widget&cssUrl=https://static.garmincdn.com/com.garmin.connect/ui/css/gauth-custom-v1.2-min.css&clientId=GarminConnect&rememberMeShown=true&rememberMeChecked=false&createAccountShown=true&openCreateAccount=false&usernameShown=false&displayNameShown=false&consumeServiceTicket=false&initialFocus=true&embedWidget=false&generateExtraServiceTicket=false")

	if err != nil {
		return err
	}

	_, err = gc.client.Get(loginUrl.String())

	if err != nil {
		return err
	}

	data := url.Values{}

	data.Set("username", username)
	data.Set("password", password)
	data.Set("_eventId", "submit")
	data.Set("embed", "true")
	data.Set("lt", "e1s1")
	data.Set("displayNameRequired", "false")

	response, err := gc.client.PostForm(loginUrl.String(), data)

	if err != nil {
		return err
	}

	if response.StatusCode != http.StatusOK {
		return errors.New("auth request not good")
	}

	body, err := ioutil.ReadAll(response.Body)

	if err != nil {
		return err
	}

	re := regexp.MustCompile(`ticket=([^"]+)`)
	ticket := re.FindAllStringSubmatch(string(body), 1)[0][1]
	response, err = gc.client.Get(GARMIN_CONNECT_URL + "/post-auth/login?ticket=" + ticket)

	if err != nil {
		return err
	}

	defer response.Body.Close()

	return nil
}
