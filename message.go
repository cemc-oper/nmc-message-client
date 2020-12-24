package nmc_message_client

type MonitorMessage struct {
	Source           string `json:"source"`
	MessageType      string `json:"type"`
	Status           string `json:"status"`
	DateTime         int64  `json:"datetime,omitempty"`
	FileName         string `json:"fileName"`
	AbsoluteDataName string `json:"absoluteDataName,omitempty"`
	Description      string `json:"desc,omitempty"`
}

type MonitorMessageV2 struct {
	Topic            string `json:"topic"`
	Source           string `json:"source"`
	SourceIP         string `json:"sourceIP"`
	MessageType      string `json:"type"`
	PID              string `json:"PID"`
	ID               string `json:"ID"`
	DateTime         string `json:"datetime,omitempty"`
	FileNames        string `json:"fileNames"`
	AbsoluteDataName string `json:"absoluteDataName,omitempty"`
	FileSizes        string `json:"fileSizes"`
	Result           int8   `json:"result"`
}
