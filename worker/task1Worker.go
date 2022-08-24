package worker

import (
	"encoding/json"

	"github.com/sirupsen/logrus"
)

// implement custom definition
type Task1Payload struct {
	Id   int    `json:"id"`
	Name string `json:"name"`
}

func PerformTask1(log *logrus.Logger, payload map[string]interface{}) error {
	// typecast payload to custom definition
	var task1Payload Task1Payload

	raw, err := json.Marshal(payload)
	if err != nil {
		log.Errorf("Error decoding JSON : %v", err.Error())
	}

	err = json.Unmarshal(raw, &task1Payload)
	if err != nil {
		log.Errorf("Error converting JSON to struct : %v", err.Error())
	}

	// implement your custom logic
	log.Infof("Performing task 1 for id : %d , name: %s", task1Payload.Id, task1Payload.Name)

	return nil
}
