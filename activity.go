package garminconnect

import (
	"bytes"
	"encoding/json"
	"io"
	"mime/multipart"
	"os"
	"path/filepath"
)

type ResultMessage struct {
	Code    int    `json:"code"`
	Content string `json:"content"`
}

type Result struct {
	InternalId int             `json:"internalId"`
	ExternalId string          `json:"externalId"`
	Messages   []ResultMessage `json:"messages"`
}

type Upload struct {
	DetailedImportResult struct {
		UploadId       int    `json:"uploadId"`
		Owner          int    `json:"owner"`
		FileSize       int    `json;"fileSize"`
		ProcessingTime int    `json:"processingTime"`
		CreationDate   string `json:"creationDate"`
		IpAddress      string `json:"ipAddress"`
		FileName       string `json:"fileName"`
		Report         struct {
			Class         string        `json:"@class"`
			CreatedOn     string        `json:"createdOn"`
			UserProfileId int           `json:"userProfileId"`
			Children      []interface{} `json:"children"`
			Entries       []interface{} `json:"entries"`
		}
		Successes []Result `json:"successes"`
		Failures  []Result `json:"failures"`
	} `json:"detailedImportResult"`
}

func (gc *Client) UploadActivity(path string) (Upload, error) {
	file, err := os.Open(path)

	if err != nil {
		return Upload{}, err
	}

	defer file.Close()

	formData := bytes.Buffer{}
	writer := multipart.NewWriter(&formData)
	activity, err := writer.CreateFormFile("data", filepath.Base(path))

	if err != nil {
		return Upload{}, err
	}
	_, err = io.Copy(activity, file)

	if err != nil {
		return Upload{}, err
	}

	writer.WriteField("responseContentType", "application/json")

	contentType := writer.FormDataContentType()
	err = writer.Close()

	if err != nil {
		return Upload{}, err
	}

	response, err := gc.client.Post(GARMIN_CONNECT_URL+"/proxy/upload-service-1.1/json/upload/.fit", contentType, &formData)

	if err != nil {
		return Upload{}, err
	}

	defer response.Body.Close()

	var upload Upload

	json.NewDecoder(response.Body).Decode(&upload)

	return upload, nil
}
