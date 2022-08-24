package consumer

import (
	"github.com/pythonista7/go-kafka-consumer/utils"
	"github.com/pythonista7/go-kafka-consumer/worker"
	"github.com/sirupsen/logrus"
)

// Main Consumer
func ConsumeQueue(
	log *logrus.Logger,
	consumeQueue <-chan utils.KafkaPayload,
) error {
	for payload := range consumeQueue {
		log.Debugf("Processing payload ...")

		switch payload.TaskType {
		case "taskType1":
			{
				err := worker.PerformTask1(log, payload.Data)
				if err != nil {
					log.Errorf("Error performing task1 : %s", err.Error())
					continue
				}
			}
			/*
				Add more cases as necessary
			*/
		default:
			{
				log.Errorf("Invalid task type: %s", payload.TaskType)
			}
		}
	}

	return nil
}
