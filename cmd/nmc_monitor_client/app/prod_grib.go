package app

import (
	"encoding/json"
	"fmt"
	"github.com/spf13/pflag"
)

func createProdGribMessage(args []string) ([]byte, error) {
	var startTime = ""
	var forecastTime = ""

	var prodGribFlagSet = pflag.NewFlagSet("prod_grid", pflag.ContinueOnError)
	prodGribFlagSet.ParseErrorsWhitelist = pflag.ParseErrorsWhitelist{UnknownFlags: true}
	prodGribFlagSet.SortFlags = false

	prodGribFlagSet.StringVar(&startTime, "start-time", "", "start time, such as 2019062400")
	prodGribFlagSet.StringVar(&forecastTime, "forecast-time", "", "forecast time, such as 000")
	if err := prodGribFlagSet.Parse(args); err != nil {
		return nil, fmt.Errorf("argument parse fail: %s", err)
	}

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
