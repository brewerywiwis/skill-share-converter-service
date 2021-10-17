package storage

import (
	"bytes"
	"errors"
	"log"
	"net/http"
	"os"
	"path/filepath"
	"skillshare/converter/config"
	"strings"
	"sync"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/credentials"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/s3"
	"github.com/aws/aws-sdk-go/service/s3/s3manager"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

func createSession() (*session.Session, error) {
	s3config := config.GetS3Config()
	return session.NewSession(&aws.Config{
		Credentials: credentials.NewStaticCredentials(s3config.S3_ACCESS_KEY_ID, s3config.S3_SECRET_KEY, ""),
		Region:      aws.String(s3config.S3_REGION)},
	)
}
func UploadFile(originalName string, mimetype string, encoding string, videoSize int, videoData bytes.Buffer) (primitive.ObjectID, *s3manager.UploadOutput, error) {
	log.Printf("Uploading file")
	session, err := createSession()
	if err != nil {
		log.Println("Cannot create session")
	}
	uploader := s3manager.NewUploader(session)
	config := config.GetS3Config()
	videoId := primitive.NewObjectID()
	result, err := uploader.Upload(&s3manager.UploadInput{
		Bucket:          aws.String(config.S3_BUCKET),
		Key:             aws.String(config.S3_RAW_VIDEO_KEY + "/" + videoId.Hex()),
		Body:            bytes.NewReader(videoData.Bytes()),
		ContentType:     aws.String(mimetype),
		ContentEncoding: aws.String(encoding),
	})
	return videoId, result, nil
}

func DeleteFile(path string) error {
	session, err := createSession()
	if err != nil {
		log.Println("Cannot create session")
		return err
	}
	svc := s3.New(session)
	config := config.GetS3Config()
	_, err = svc.DeleteObject(&s3.DeleteObjectInput{
		Bucket: aws.String(config.S3_BUCKET),
		Key:    aws.String(path),
	})
	if err != nil {
		log.Println("Cannot create session")
		return err
	}
	err = svc.WaitUntilObjectNotExists(&s3.HeadObjectInput{
		Bucket: aws.String(config.S3_BUCKET),
		Key:    aws.String(path),
	})
	if err != nil {
		return err
	}
	return nil
}

func UploadDirToS3(dir string) {
	fileList := []string{}
	filepath.Walk(dir, func(path string, f os.FileInfo, err error) error {
		fileList = append(fileList, path)
		return nil
	})
	var wg sync.WaitGroup
	wg.Add(len(fileList) - 1)
	for _, pathOfFile := range fileList[1:] {
		go putInS3(pathOfFile, &wg)
	}
	wg.Wait()
	log.Println("All files are uploaded")
}
func PreprocessPath(path string, tmpDirName string) (string, error) {
	i := strings.LastIndex(path, tmpDirName)
	if i == -1 {
		return "", errors.New("cannot find tmp dir")
	}
	return strings.Trim(path[i+len(tmpDirName):], "/"), nil
}
func putInS3(pathOfFile string, wg *sync.WaitGroup) {
	defer func() {
		wg.Done()
	}()
	file, _ := os.Open(pathOfFile)
	defer file.Close()
	fileInfo, _ := file.Stat()
	size := fileInfo.Size()
	buffer := make([]byte, size)
	file.Read(buffer)
	fileBytes := bytes.NewReader(buffer)
	fileType := http.DetectContentType(buffer)
	path := file.Name()
	tmpDirName := "tmp"
	fileName, err := PreprocessPath(path, tmpDirName)
	if err != nil {
		log.Println("Cannot find filename")
		return
	}
	session, err := createSession()
	if err != nil {
		log.Println("Cannot create session")
	}
	uploader := s3manager.NewUploader(session)
	config := config.GetS3Config()
	_, err = uploader.Upload(&s3manager.UploadInput{
		Bucket:      aws.String(config.S3_BUCKET),
		Key:         aws.String(config.S3_HLS_VIDEO_KEY + "/" + fileName),
		Body:        fileBytes,
		ContentType: aws.String(fileType),
	})
}

func DownloadFromS3Bucket(key, pathToSave string) error {
	file, err := os.Create(pathToSave)
	if err != nil {
		log.Println("Cannot create file with path: ", err)
		return err
	}
	defer file.Close()

	session, err := createSession()
	if err != nil {
		log.Println("Cannot create session")
		return err
	}
	downloader := s3manager.NewDownloader(session)
	config := config.GetS3Config()
	_, err = downloader.Download(file,
		&s3.GetObjectInput{
			Bucket: aws.String(config.S3_BUCKET),
			Key:    aws.String(key),
		})
	if err != nil {
		log.Println("Cannot download object with key")
		return err
	}

	log.Println("Downloaded")
	return nil
}
