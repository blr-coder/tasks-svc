package config

import "github.com/ilyakaznacheev/cleanenv"

const appConfigPath = "configs/appconfig.yaml"

type AppConfig struct {
	PostgresConnLink string      `yaml:"postgres_conn_link" env:"POSTGRES_CONN_LINK" env-required:"true"`
	KafkaConfig      KafkaConfig `yaml:"kafka_config"`
	AppPort          string      `yaml:"app_port"`
}

type KafkaConfig struct {
	Topic     string `yaml:"topic" env:"TOPIC" env-default:"async_arc_topic"`
	Partition int    `yaml:"partition" env:"PARTITION" env-default:"0"`
	Address   string `yaml:"address" env:"ADDRESS" env-default:"localhost:9092"`
}

func NewAppConfig() (*AppConfig, error) {
	var cfg AppConfig

	err := cleanenv.ReadConfig(appConfigPath, &cfg)
	if err != nil {
		return nil, err
	}

	return &cfg, nil
}
