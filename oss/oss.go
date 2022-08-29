package oss

import (
	"bytes"
	"fmt"

	"github.com/aliyun/aliyun-oss-go-sdk/oss"
	"github.com/rs/zerolog/log"
)

type OSSCOnfig struct {
	Endpoint        string
	AccessKeyId     string
	AccessKeySecret string
	Bucket          string
	Url             string
}

type OSSInterface interface {
	Start() (client *oss.Client, err error)
	Upload(key string, object []byte) (url string, err error)
}

type ossInstance struct {
	*OSSCOnfig
	client *oss.Client
}

func NewClient(conf *OSSCOnfig) OSSInterface {
	return &ossInstance{
		conf,
		nil,
	}
}

func (o *ossInstance) Start() (client *oss.Client, err error) {
	client, err = oss.New(o.Endpoint, o.AccessKeyId, o.AccessKeySecret)
	if err != nil {
		log.Error().Err(err).Msg("error on strating oss client")
		return
	}

	// create bucket
	exist, err := client.IsBucketExist(o.Bucket)
	if err != nil {
		log.Error().Err(err).Msg("error on check bucket")
		return
	}
	if !exist {
		err = client.CreateBucket(o.Bucket)
		if err != nil {
			log.Error().Err(err).Msg("error on create bucket")
			return
		}
	}

	o.client = client
	return
}

func (o *ossInstance) Upload(key string, object []byte) (url string, err error) {
	bucket, err := o.client.Bucket(o.Bucket)
	if err != nil {
		log.Error().Err(err).Msg("error on access bucket")
		return
	}

	err = bucket.PutObject(key, bytes.NewReader(object), oss.ObjectACL(oss.ACLPublicRead))
	if err != nil {
		log.Error().Err(err).Msg("error on put file to bucket")
		return
	}
	url = fmt.Sprintf(o.Url+"/%s", key)
	return
}
