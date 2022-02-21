package settings

import (
	"context"
	"os"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/feature/s3/manager"
	"github.com/aws/aws-sdk-go-v2/service/s3"
	"github.com/sirupsen/logrus"
	"gopkg.in/yaml.v3"
)

func NewS3Config(cfg aws.Config, key string, bucket string) (Settings, error) {
	var buf []byte
	s3conf := manager.NewWriteAtBuffer(buf)
	s3Client := s3.NewFromConfig(cfg)

	downloader := manager.NewDownloader(s3Client)
	if _, err := downloader.Download(context.TODO(), s3conf, &s3.GetObjectInput{
		Bucket: aws.String(bucket),
		Key:    aws.String(key),
	}); err != nil {
		logrus.Fatal(err.Error())
	}

	return loadConfigFromBytes(s3conf.Bytes())
}

func NewLocalConfig(path string) (Settings, error) {
	var Conf Settings

	resp, err := os.ReadFile(path)
	if err != nil {
		return Settings{}, err
	}

	if err = yaml.Unmarshal(resp, &Conf); err != nil {
		return Settings{}, err
	}

	return Conf, nil
}
