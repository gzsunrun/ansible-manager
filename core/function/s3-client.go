package function

import (
	"io"
	"net/url"
	"time"
	"strings"
	"net"

	log "github.com/astaxie/beego/logs"
	"github.com/gzsunrun/ansible-manager/core/config"
	"github.com/minio/minio-go"
)

type s3Client struct {
	URL    string
	Key    string
	Secret string
	Bucket string
}

var s3 = s3Client{}

func NewS3Client() {
	s3.URL = config.Cfg.Ansible.S3URL
	s3.Key = config.Cfg.Ansible.S3Key
	s3.Secret = config.Cfg.Ansible.S3Secret
	s3.Bucket = config.Cfg.Ansible.BucketName
}

func S3Put(reader io.Reader, size int64, fileName string) error {
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

func GeneratePresignedUrl(object string) (string, error) {
	addrs:=strings.Split(s3.URL,":")
	ns, err := net.LookupHost(addrs[0])
	if err!=nil{
		return "",err
	}
	s3Client, err := minio.NewV2(ns[0]+":"+addrs[1], s3.Key, s3.Secret, false)
	if err != nil {
		log.Error(err)
		return "", err
	}
	object += ".png"
	reqParams := make(url.Values)
	reqParams.Set("response-content-disposition", "attachment;filename=\""+object+"\"")

	presignedURL, err := s3Client.PresignedGetObject(s3.Bucket, object, time.Second*60*10, reqParams)
	if err != nil {
		log.Error(err)
		return "", err
	}
	return presignedURL.String(), nil
}
