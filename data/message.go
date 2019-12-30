package data

type Message struct {
	Timestamp int64       `json:"timestamp"`
	Value     interface{} `json:"value"`
}
