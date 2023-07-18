package api

type ResponseV1 struct {
	Status bool `json:"status"`
	Error  error
	Data   ResponseV1Data `json:"data"`
}

type ResponseV1Data struct {
	Time string `json:"time"`
}
