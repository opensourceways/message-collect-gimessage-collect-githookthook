package config

import "github.com/opensourceways/message-collect-githook/kafka"

type Config struct {
	Kafka kafka.Config `json:"kafka"`
	Port  int          `yaml:"port"`
}
