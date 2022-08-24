package utils

type KafkaPayload struct {
	TaskType string                 `json:"taskType"`
	Data     map[string]interface{} `json:"data"`
}
