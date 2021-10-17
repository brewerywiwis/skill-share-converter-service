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
	URL     string
	DB_NAME string
}

type S3Config struct {
	S3_REGION        string
	S3_BUCKET        string
	S3_ACCESS_KEY_ID string
	S3_SECRET_KEY    string
	S3_RAW_VIDEO_KEY string
	S3_HLS_VIDEO_KEY string
}

var rabbitMQ *RabbitConfig
var databaseConfig *DatabaseConfig
var s3Config *S3Config

func Init() {
	// viper.AutomaticEnv()
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
			URL:     viper.GetString("mongo.url"),
			DB_NAME: viper.GetString("mongo.db_name"),
		}
	}
	return databaseConfig
}

func GetS3Config() *S3Config {
	if s3Config == nil {
		s3Config = &S3Config{
			S3_REGION:        viper.GetString("aws.s3_region"),
			S3_BUCKET:        viper.GetString("aws.s3_bucket"),
			S3_ACCESS_KEY_ID: viper.GetString("aws.s3_access_key_id"),
			S3_SECRET_KEY:    viper.GetString("aws.s3_secret_key"),
			S3_RAW_VIDEO_KEY: viper.GetString("aws.s3_raw_video_key"),
			S3_HLS_VIDEO_KEY: viper.GetString("aws.s3_hls_video_key"),
		}
	}
	return s3Config
}
