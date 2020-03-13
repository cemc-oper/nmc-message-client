package nmc_message_client

import (
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
	Offset       string    `json:"_id"`
	Source       string    `json:"source"`
	MessageType  string    `json:"type"`
	Status       string    `json:"status"`
	DateTime     time.Time `json:"datetime,omitempty"`
	FileName     string    `json:"fileName"`
	StartTime    string    `json:"startTime,omitempty"`
	ForecastTime string    `json:"forecastTime,omitempty"`
}
