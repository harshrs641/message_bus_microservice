package event

type Event struct {
	ServiceName string `json:"serviceName"`
	Data        string `json:"data"`
}
