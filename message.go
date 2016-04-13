package garminconnect

import (
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
)

const (
	WORKOUT_FILE_TYPE = "FIT_TYPE_5"
)

type Queue struct {
	Host             string    `json:"serviceHost"`
	NumberOfMessages int       `json:"numOfMessages"`
	Messages         []Message `json:"messages"`
}

type Message struct {
	Id                int      `json:"messageId"`
	Type              string   `json:"messageType"`
	Status            string   `json:"messageStatus"`
	DeviceId          int      `json:"deviceId"`
	DeviceName        string   `json:"deviceName"`
	ApplicationKey    string   `json:"applicationKey"`
	FirmwareVersion   string   `json:"FirmwareVersion"`
	WifiSetup         bool     `json:"wifiSetup"`
	DeviceXmlDataType string   `json:"deviceXmlDataType"`
	Metadata          Metadata `json:"metadata"`
}

type Metadata struct {
	Filetype    string `json:"fileType"`
	MessageUrl  string `json:"messageUrl"`
	Absolute    bool   `json:"absolute"`
	MessageName string `json:"messageName"`
	GroupName   string `json:"groupName"`
	Priority    int    `json:"priority"`
	Id          int    `json:"metaDataId"`
	AppDetails  string `json:"appDetails"`
}

func (gc *Client) Messages() ([]Message, error) {
	response, err := gc.client.Get(GARMIN_CONNECT_URL + "/modern/proxy/device-service/devicemessage/messages")

	if err != nil {
		return []Message{}, err
	}

	defer response.Body.Close()

	var queue Queue

	err = json.NewDecoder(response.Body).Decode(&queue)

	if err != nil {
		panic(err)
	}

	return queue.Messages, nil
}

func (gc *Client) DeleteMessage(messageId int) error {
	request, err := http.NewRequest(http.MethodDelete, GARMIN_CONNECT_URL+fmt.Sprintf("/modern/proxy/device-service/devicemessage/message/%d", messageId), nil)

	if err != nil {
		return err
	}

	response, err := gc.client.Do(request)

	if err != nil {
		return err
	}

	defer response.Body.Close()

	if response.StatusCode != http.StatusOK {
		return errors.New(fmt.Sprintf("%d", response.StatusCode))
	}

	return nil
}
