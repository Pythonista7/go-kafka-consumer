package main

import (
	"encoding/json"
	"fmt"
	"os"
	"os/signal"
	"strings"
	"syscall"
	"time"

	"github.com/pythonista7/go-kafka-consumer/config"
	"github.com/pythonista7/go-kafka-consumer/consumer"
	"github.com/pythonista7/go-kafka-consumer/utils"
	"github.com/sirupsen/logrus"
	"gopkg.in/confluentinc/confluent-kafka-go.v1/kafka"
)

func main() {

	log := logrus.New()
	log.Out = os.Stdout
	log.SetReportCaller(true)

	if config.Config.AppEnv != "production" {
		log.SetLevel(logrus.TraceLevel)
	}

	if config.Config.LogFormat == "json" {
		log.SetFormatter(&logrus.JSONFormatter{})
	}

	topics := strings.Split(config.Config.KafkaTopic, ",")
	fmt.Printf("Broker : %s\n", config.Config.KafkaBroker)
	sigchan := make(chan os.Signal, 1)
	signal.Notify(sigchan, syscall.SIGINT, syscall.SIGTERM)

	c, err := kafka.NewConsumer(&kafka.ConfigMap{
		"allow.auto.create.topics":        true,
		"bootstrap.servers":               config.Config.KafkaBroker,
		"group.id":                        config.Config.KafkaGroup,
		"session.timeout.ms":              6000,
		"go.events.channel.enable":        true,
		"go.application.rebalance.enable": true,
		// Enable generation of PartitionEOF when the
		// end of a partition is reached.
		"enable.partition.eof": true,
		"auto.offset.reset":    "earliest"})

	if err != nil {
		fmt.Fprintf(os.Stderr, "Failed to create consumer: %s\n", err)
		os.Exit(1)
	}

	log.Infof("Created Consumer %s\n", c)

	//define channel
	consume := make(chan utils.KafkaPayload, config.Config.ChannelBufferSize)

	// spawn go routine in block state for the channel above
	for i := 0; i < config.Config.NoOfThreads; i++ {
		go consumer.ConsumeQueue(log, consume)
	}

	// Note: you can also use the non-channel based approach, using the poll API as shown here:
	// https://github.com/confluentinc/confluent-kafka-go/blob/master/examples/consumer_example/consumer_example.go

	err = c.SubscribeTopics(topics, nil)
	if err != nil {
		log.Errorf("Error while subscribing topics: %v | Error: %s", topics, err.Error())
	}

	run := true

	for run == true {

		select {
		case sig := <-sigchan:
			log.Infof("Caught signal %s: terminating\n", sig.String())
			run = false

		case ev := <-c.Events():
			switch e := ev.(type) {
			case kafka.AssignedPartitions:
				log.Infof("AssignedPartitions : %v\n", e)
				fmt.Fprintf(os.Stderr, "%% %v\n", e)
				c.Assign(e.Partitions)

			case kafka.RevokedPartitions:
				log.Errorf("RevokedPartitions : %v\n", e)
				fmt.Fprintf(os.Stderr, "%% %v\n", e)
				c.Unassign()

			case *kafka.Message:
				log.Debugf("Buffer occupency %d / %d", len(consume), config.Config.ChannelBufferSize)
				log.Debugf("Picked up kafka message at : %v", time.Now())
				log.Infoln("Processing Kafka Message ...")

				var payload utils.KafkaPayload
				err = json.Unmarshal(e.Value, &payload)
				consume <- payload

			case kafka.PartitionEOF:
				log.Infof("Reached %v\n", e)

			case kafka.Error:
				// Errors should generally be considered as informational, the client will try to automatically recover
				fmt.Fprintf(os.Stderr, "%% Error: %v\n", e)
				log.Errorf("Error: %v\n", e)
			}
		}
	}
	fmt.Printf("Closing consumer\n")
	c.Close()
}
