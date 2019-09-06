package app

type MonitorMessage struct {
	Source           string `json:"source"`
	MessageType      string `json:"type"`
	Status           string `json:"status"`
	DateTime         int64  `json:"datetime,omitempty"`
	FileName         string `json:"fileName"`
	AbsoluteDataName string `json:"absoluteDataName,omitempty"`
	Description      string `json:"desc,omitempty"`
}
