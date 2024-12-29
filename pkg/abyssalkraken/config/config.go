package config

import (
	"github.com/spf13/viper"
)

type AbyssalKrakenProperties struct {
	Events    EventsConfig    `mapstructure:"events"`
	Snapshots SnapshotsConfig `mapstructure:"snapshots"`
}

type EventsConfig struct {
	Persistence   PersistenceConfig   `mapstructure:"persistence"`
	Serialization SerializationConfig `mapstructure:"serialization"`
}

type SnapshotsConfig struct {
	Persistence   PersistenceConfig   `mapstructure:"persistence"`
	Serialization SerializationConfig `mapstructure:"serialization"`
	Thresholds    map[string]int      `mapstructure:"thresholds"`
}

type PersistenceConfig struct {
	R2dbc R2dbcPersistenceConfig `mapstructure:"r2dbc"`
}

type R2dbcPersistenceConfig struct {
	TableName            string `mapstructure:"tableName"`
	ValidateSchemaOnInit bool   `mapstructure:"validateSchemaOnInit" default:"true"`
}

type SerializationConfig struct {
	Avro AvroSerializationConfig `mapstructure:"avro"`
}

type AvroSerializationConfig struct {
	Encoding string `mapstructure:"encoding"`
}

func LoadConfig(path string) (*AbyssalKrakenProperties, error) {
	viper.SetConfigFile(path + "/abyssal_kraken.yaml")
	viper.SetConfigType("yaml")
	viper.AddConfigPath(path)
	viper.AutomaticEnv()

	var config AbyssalKrakenProperties

	if err := viper.ReadInConfig(); err != nil {
		return nil, err
	}

	if err := viper.Unmarshal(&config); err != nil {
		return nil, err
	}

	return &config, nil
}
