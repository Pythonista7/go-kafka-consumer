package config

import (
	"github.com/ilyakaznacheev/cleanenv"
	"github.com/sirupsen/logrus"
)

type Configuration struct {
	AppEnv            string `env:"APP_ENV" env-description:"app environment" env-default:"staging"`
	LogFormat         string `env:"LOG_FORMAT" env-description:"format of the go-consumer log" env-default:"text"`
	NoOfThreads       int    `env:"NO_OF_THREADS" env-description:"Total number of threads/go-routines to be spawned" env-default:"4"`
	ChannelBufferSize int    `env:"CHANNEL_BUFFER_SIZE" env-description:"total size of channel buffer" env-default:"5"`

	KafkaBroker string `env:"KAFKA_BROKER" env-description:"kafka broker address host:port" env-default:"localhost:9092"`
	KafkaGroup  string `env:"KAFKA_GROUP" env-description:"kafka consumer group" env-default:"consumer-0"`
	KafkaTopic  string `env:"KAFKA_TOPIC" env-description:"kafka topic" env-default:"test-topic-1"`
}

var Config Configuration

func init() {
	err := cleanenv.ReadEnv(&Config)
	if err != nil {
		logrus.Errorf("Something went wrong while loading configurations from env | Error: %s", err.Error())
	}
}
