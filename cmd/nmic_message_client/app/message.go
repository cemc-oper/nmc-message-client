package app

type MonitorMessage struct {
	AtMessageType string `json:"@type"`
	AtName        string `json:"@name"`
	AtMessage     string `json:"@message"`
	AtOccurTime   int64  `json:"@occur_time"`
	DataType      string `json:"DATA_TYPE"`
	DataType1     string `json:"DATA_TYPE_1"`
	DataTime      string `json:"DATA_TIME"`
	FileNameO     string `json:"FILE_NAME_O"`
	Send          string `json:"SEND"`
	Receive       string `json:"RECEIVE"`
	BusinessState string `json:"BUSINESS_STATE"`
}
