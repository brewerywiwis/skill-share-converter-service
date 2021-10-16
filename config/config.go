package config

import (
	"fmt"

	"github.com/spf13/viper"
)

type RabbitConfig struct {
	USERNAME              string
	PASSWORD              string
	HOST                  string
	PORT                  string
	SensorGatewayExchange string
	RoutingKeySuffix      string
}

type DatabaseConfig struct {
	URL        string
	DB_NAME    string
	BACKEND_DB string
}

var rabbitMQ *RabbitConfig
var databaseConfig *DatabaseConfig

func Init() {
	viper.SetConfigName("config") // name of config file (without extension)
	viper.SetConfigType("yaml")   // REQUIRED if the config file does not have the extension in the name
	viper.AddConfigPath(".")      // optionally look for config in the working directory
	viper.AddConfigPath("..")     // optionally look for config in the working directory
	err := viper.ReadInConfig()   // Find and read the config file
	if err != nil {               // Handle errors reading the config file
		panic(fmt.Errorf("Fatal error config file: %w", err))
	}
}

func GetRabbitMQConfig() *RabbitConfig {
	if rabbitMQ == nil {
		rabbitMQ = &RabbitConfig{
			USERNAME:              viper.GetString("rabbit_mq.username"),
			PASSWORD:              viper.GetString("rabbit_mq.password"),
			HOST:                  viper.GetString("rabbit_mq.host"),
			PORT:                  viper.GetString("rabbit_mq.port"),
			SensorGatewayExchange: viper.GetString("rabbit_mq.sensorGatewayExchange"),
			RoutingKeySuffix:      viper.GetString("rabbit_mq.routingSuffix"),
		}

	}
	return rabbitMQ
}

func GetDatabaseConfig() *DatabaseConfig {
	if databaseConfig == nil {
		databaseConfig = &DatabaseConfig{
			URL:        viper.GetString("mongo.url"),
			DB_NAME:    viper.GetString("mongo.db_name"),
			BACKEND_DB: viper.GetString("mongo.backend_db"),
		}
	}
	return databaseConfig
}
