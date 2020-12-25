package nmc_message_client

import (
	"encoding/json"
	"fmt"
	"math/rand"
	"strconv"
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
	dateTime time.Time,
	fileName string,
	absoluteDataName string,
	fileSize int64,
	productIntervalHour int32,
	status int8,
	startTime string,
	forecastTime string,
) ([]byte, error) {
	timeString := dateTime.Format("2006-01-02 15:04:05")

	startTimeClock, _ := time.Parse("2006010215", startTime)
	var cstZone = time.FixedZone("CST", 8*3600) // beijing
	startTimeClock = startTimeClock.In(cstZone)
	forecastHour, _ := strconv.Atoi(forecastTime)
	randomNumber := rand.Intn(10000)

	description := ProbGribMessageDescription{
		StartTime:    startTime,
		ForecastTime: forecastTime,
	}
	descriptionBlob, err := json.Marshal(description)
	if err != nil {
		return nil, fmt.Errorf("create description for prod-grib error: %s", err)
	}

	monitorMessage := MonitorMessageV2{
		Topic:             topic,
		Source:            source,
		SourceIP:          sourceIP,
		MessageType:       messageType,
		Result:            status,
		DateTime:          timeString,
		FileNames:         fileName,
		AbsoluteDataName:  absoluteDataName,
		FileSizes:         fmt.Sprintf("%d", fileSize),
		ResultDescription: string(descriptionBlob),
	}

	monitorMessage.ID = fmt.Sprintf(
		"%s%05d",
		dateTime.Format("20060102150405"),
		dateTime.Nanosecond()/10000,
	)
	monitorMessage.PID = fmt.Sprintf(
		"%s00%02d%04d%04d%04d",
		messageType,
		startTimeClock.Hour(),
		forecastHour,
		productIntervalHour,
		randomNumber,
	)

	monitorMessageBlob, err := json.MarshalIndent(monitorMessage, "", "  ")

	if err != nil {
		return monitorMessageBlob, fmt.Errorf("create message for prod-grib error: %s", err)
	}

	return monitorMessageBlob, nil
}
