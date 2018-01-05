package s3

import (
	"io"

	"git.gzsunrun.cn/sunrunlib/log"
	"github.com/gzsunrun/ansible-manager/config"
	"github.com/minio/minio-go"
)

type s3Client struct {
	URL    string
	Key    string
	Secret string
	Bucket string
	Status bool
}

var s3 = s3Client{}

func NewClient() {
	s3.URL = config.Cfg.AnsibleManager.S3URL
	s3.Key = config.Cfg.AnsibleManager.S3Key
	s3.Secret = config.Cfg.AnsibleManager.S3Secret
	s3.Bucket = config.Cfg.AnsibleManager.BucketName
	s3.Status = config.Cfg.AnsibleManager.S3Status
}

func S3Put(reader io.Reader, size int64, fileName string) error {
	if !s3.Status {
		return nil
	}
	s3Client, err := minio.NewV2(s3.URL, s3.Key, s3.Secret, false)
	if err != nil {
		log.Error(err)
		return err
	}

	_, err = s3Client.PutObject(s3.Bucket, fileName, reader, size, minio.PutObjectOptions{ContentType: "application/octet-stream"})
	if err != nil {
		log.Error(err)
		return err
	}
	return nil
}

func S3Get(object string, savePath string) error {
	if !s3.Status {
		return LocalCopy(object, savePath)
	}
	s3Client, err := minio.NewV2(s3.URL, s3.Key, s3.Secret, false)
	if err != nil {
		log.Error(err)
		return err
	}

	err = s3Client.FGetObject(s3.Bucket, object, savePath, minio.GetObjectOptions{})
	if err != nil {
		log.Error(err)
		return err
	}
	return nil
}

func S3Delte(object string) error {
	if !s3.Status {
		return LocalDel(object)
	}
	s3Client, err := minio.NewV2(s3.URL, s3.Key, s3.Secret, false)
	if err != nil {
		log.Error(err)
		return err
	}
	err = s3Client.RemoveObject(s3.Bucket, object)
	if err != nil {
		log.Error(err)
		return err
	}
	return nil
}
