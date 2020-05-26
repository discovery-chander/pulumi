package s3

import (
	"github.com/EurosportDigital/global-transcoding-platform/lib/errors"
	"github.com/aws/aws-sdk-go/aws/awserr"
	awsS3 "github.com/aws/aws-sdk-go/service/s3"
	"github.com/aws/aws-sdk-go/service/s3/s3iface"
)

type S3Client interface {
	CheckResourceExistance(bucket string, key string) (bool, error)
}

type s3ClientObject struct {
	awsS3Client s3iface.S3API
}

const (
	resourceNotFoundCode = "NotFound"
)

func New(awsS3Client s3iface.S3API) *s3ClientObject {
	return &s3ClientObject{
		awsS3Client: awsS3Client,
	}
}

func (instance *s3ClientObject) CheckResourceExistance(bucket string, key string) (bool, error) {
	_, err := instance.awsS3Client.HeadObject(&awsS3.HeadObjectInput{
		Bucket: &bucket,
		Key:    &key,
	})
	if err != nil {
		if aerr, ok := err.(awserr.Error); ok && aerr.Code() == resourceNotFoundCode {
			return false, nil
		}
		return false, errors.Wrap(err, "an unexpected error ocurred while finding source input")
	}
	return true, nil
}
