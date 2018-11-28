package storage

import (
	"io/ioutil"
	"net"
	"net/url"
	"path"
	"strings"
	"time"

	"github.com/hashwing/log"
	"github.com/minio/minio-go"
)

// S3Storage s3 repo
type S3Storage struct {
	URL    string
	Key    string
	Secret string
	Bucket string
}

// NewS3Storage new s3 storage
func NewS3Storage(url, key, secret, bucket string) (*S3Storage, error) {
	s3 := new(S3Storage)
	s3.URL = url
	s3.Key = key
	s3.Secret = secret
	s3.Bucket = bucket
	err := s3.CreateBuket(bucket)
	return s3, err
}

// CreateBuket create s3 buket
func (s3 *S3Storage) CreateBuket(name string) error {
	s3Client, err := minio.NewV2(s3.URL, s3.Key, s3.Secret, false)
	if err != nil {
		log.Error(err)
		return err
	}
	res, err := s3Client.BucketExists(name)
	if err != nil {
		log.Error(err)
		return err
	}
	if res {
		return nil
	}
	return s3Client.MakeBucket(name, "us-east-1")
}

// Put upload file
func (s3 *S3Storage) Put(repo *StorageParse) error {
	s3Client, err := minio.NewV2(s3.URL, s3.Key, s3.Secret, false)
	if err != nil {
		log.Error(err)
		return err
	}
	_, err = s3Client.FPutObject(s3.Bucket, repo.RemotePath, repo.LocalPath, minio.PutObjectOptions{})
	if err != nil {
		log.Error(err)
		return err
	}
	return nil
}

// Get download file
func (s3 *S3Storage) Get(repo *StorageParse) error {
	s3Client, err := minio.NewV2(s3.URL, s3.Key, s3.Secret, false)
	if err != nil {
		log.Error(err)
		return err
	}

	err = s3Client.FGetObject(s3.Bucket, repo.RemotePath, repo.LocalPath, minio.GetObjectOptions{})
	if err != nil {
		log.Error(err)
		return err
	}
	return nil
}

// Delete delete file
func (s3 *S3Storage) Delete(repo *StorageParse) error {
	s3Client, err := minio.NewV2(s3.URL, s3.Key, s3.Secret, false)
	if err != nil {
		log.Error(err)
		return err
	}
	err = s3Client.RemoveObject(s3.Bucket, repo.RemotePath)
	if err != nil {
		log.Error(err)
		return err
	}
	return nil
}

// Share share file
func (s3 *S3Storage) Share(repo *StorageParse) (string, error) {
	addrs := strings.Split(s3.URL, ":")
	ns, err := net.LookupHost(addrs[0])
	if err != nil {
		return "", err
	}
	s3Client, err := minio.NewV2(ns[0]+":"+addrs[1], s3.Key, s3.Secret, false)
	if err != nil {
		log.Error(err)
		return "", err
	}
	reqParams := make(url.Values)
	reqParams.Set("response-content-disposition", "attachment;filename=\""+repo.RemotePath+"\"")

	presignedURL, err := s3Client.PresignedGetObject(s3.Bucket, repo.RemotePath, time.Second*60*10, reqParams)
	if err != nil {
		log.Error(err)
		return "", err
	}
	return presignedURL.String(), nil
}

func (s3 *S3Storage) GetIO(repo *StorageParse) ([]byte, string, error) {
	s3Client, err := minio.NewV2(s3.URL, s3.Key, s3.Secret, false)
	if err != nil {
		log.Error(err)
		return nil, "", err
	}
	doneCh := make(chan struct{})
	defer close(doneCh)
	for object := range s3Client.ListObjects(s3.Bucket, "logos/"+repo.RemotePath, true, doneCh) {
		if object.Err != nil {
			return nil, "", object.Err
		}
		o, err := s3Client.GetObject(s3.Bucket, object.Key, minio.GetObjectOptions{})
		if err != nil {
			log.Error(err)
			return nil, "", err
		}
		defer o.Close()
		data, err := ioutil.ReadAll(o)
		return data, path.Ext(object.Key), err
	}
	return nil, "", nil
}
