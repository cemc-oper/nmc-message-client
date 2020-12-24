package nmc_message_client

import (
	"crypto/md5"
	"encoding/json"
	"fmt"
	"time"
)

func CreateProdGribMessage(
	source string,
	messageType string,
	status string,
	datetime int64,
	fileName string,
	absoluteDataName string,
	startTime string,
	forecastTime string) ([]byte, error) {

	description := ProbGribMessageDescription{
		StartTime:    startTime,
		ForecastTime: forecastTime,
	}
	descriptionBlob, err := json.Marshal(description)
	if err != nil {
		return nil, fmt.Errorf("create description for prod-grib error: %s", err)
	}

	monitorMessage := MonitorMessage{
		Source:           source,
		MessageType:      messageType,
		Status:           status,
		DateTime:         datetime,
		FileName:         fileName,
		AbsoluteDataName: absoluteDataName,
		Description:      string(descriptionBlob),
	}

	monitorMessageBlob, err := json.MarshalIndent(monitorMessage, "", "  ")

	if err != nil {
		return monitorMessageBlob, fmt.Errorf("create message for prod-grib error: %s", err)
	}

	return monitorMessageBlob, nil
}

type ProbGribMessageDescription struct {
	StartTime    string `json:"startTime,omitempty"`
	ForecastTime string `json:"forecastTime,omitempty"`
}

type GribProduction struct {
	Offset       string    `json:"-"`
	Source       string    `json:"source"`
	MessageType  string    `json:"type"`
	Status       string    `json:"status"`
	DateTime     time.Time `json:"datetime,omitempty"`
	FileName     string    `json:"fileName"`
	StartTime    time.Time `json:"startTime,omitempty"`
	ForecastTime string    `json:"forecastTime,omitempty"`
}

func CreateProdGribMessageV2(
	topic string,
	source string,
	sourceIP string,
	messageType string,
	dateTime string,
	fileName string,
	absoluteDataName string,
	fileSize int64,
	status int8,
	startTime string,
	forecastTime string,
) ([]byte, error) {
	monitorMessage := MonitorMessageV2{
		Topic:            topic,
		Source:           source,
		SourceIP:         sourceIP,
		MessageType:      messageType,
		Result:           status,
		DateTime:         dateTime,
		FileNames:        fileName,
		AbsoluteDataName: absoluteDataName,
		FileSizes:        fmt.Sprintf("%d", fileSize),
	}

	monitorMessage.ID = fmt.Sprintf("%x", md5.Sum([]byte(dateTime)))
	monitorMessage.PID = fmt.Sprintf(
		"%s+%s+%s+0000",
		messageType,
		startTime,
		forecastTime,
	)

	monitorMessageBlob, err := json.MarshalIndent(monitorMessage, "", "  ")

	if err != nil {
		return monitorMessageBlob, fmt.Errorf("create message for prod-grib error: %s", err)
	}

	return monitorMessageBlob, nil
}
