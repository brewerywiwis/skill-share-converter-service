package main

import (
	"encoding/json"
	"log"
	"os"
	"os/exec"
	"path/filepath"
	"skillshare/converter/config"
	"skillshare/converter/mq"
	"skillshare/converter/storage"
	"strings"
)

type ConverterMessage struct {
	VideoLink string `json:"video_link"`
}

func failOnError(err error, msg string) {
	if err != nil {
		log.Fatalf("%s: %s", msg, err)
	}
}

func Convert(filePath string) {
	cmd := exec.Command("bash", "create-vod-hls.sh", filePath)
	bytes, err := cmd.Output()
	if err != nil {
		log.Println(string(bytes), err.Error())
		return
	}
	log.Println("Converted")
	directoryPath := strings.TrimSuffix(filePath, ".mp4")
	storage.UploadDirToS3(directoryPath)
	err = os.RemoveAll(directoryPath)
	if err != nil {
		log.Println(err)
	}
	log.Println("Finished")
}

func main() {
	config.Init()
	// database.Init()
	// defer database.Disconnect()
	currentQueue := "converter"
	// if len(os.Args) < 2 {
	// 	log.Printf("Usage: %s [info] [warning] [error]", os.Args[0])
	// 	os.Exit(0)
	// } else if len(os.Args) == 2 {
	// 	currentQueue = os.Args[1]
	// }

	rabbitMQ := config.GetRabbitMQConfig()

	routingKey := currentQueue + rabbitMQ.RoutingKeySuffix

	mq.CreateChannel(currentQueue, routingKey)
	defer mq.ClearConnection()

	msgs := mq.Consume()

	forever := make(chan bool)

	go func() {
		for d := range msgs {
			var data ConverterMessage
			json.Unmarshal(d.Body, &data)
			log.Printf("Received a message: %s", data)
			key, err := storage.PreprocessPath(data.VideoLink, ".com")
			if err != nil {
				log.Println("Cannot find key from message")
				return
			}
			// log.Println(key)
			root, err := os.Getwd()
			if err != nil {
				log.Println("Cant get root dir")
			}
			fileName, err := storage.PreprocessPath(key, "/")
			fileName = strings.Trim(fileName, "/") + ".mp4"
			filePath := filepath.Join(root, "tmp", fileName)
			storage.DownloadFromS3Bucket(key, filePath)
			Convert(filePath)
			err = os.Remove(filePath)
			if err != nil {
				log.Println(err)
			}
		}
	}()

	log.Printf(" [*] Waiting for messages. To exit press CTRL+C")
	<-forever
}
