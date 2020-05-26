package s3

import (
	"fmt"
	"testing"

	"github.com/EurosportDigital/global-transcoding-platform/lib/aws/s3/mocks"
	"github.com/EurosportDigital/global-transcoding-platform/lib/errors"
	"github.com/aws/aws-sdk-go/aws/awserr"
	"github.com/aws/aws-sdk-go/service/s3"
	"github.com/stretchr/testify/suite"
)

type S3TestSuite struct {
	suite.Suite
	s3Client  S3Client
	awsS3Mock *mocks.S3API
}

func TestS3TestSuite(t *testing.T) {
	suite.Run(t, new(S3TestSuite))
}

func (suite *S3TestSuite) SetupTest() {
	mock := &mocks.S3API{}
	suite.s3Client = New(mock)
	suite.awsS3Mock = mock
}

func (suite *S3TestSuite) TestAssetValidation() {
	var (
		mockedNotFoundError   = awserr.New("NotFound", "mock", fmt.Errorf("mock"))
		mockedUnexpectedError = awserr.New("mockedUnexpected", "mock", fmt.Errorf("mock"))
		validBucket           = "mockedValidBucket"
		validKey              = "mockedValidKey"
		invalidBucket         = "mockedInvalidBucket"
		invalidKey            = "mockedInvalidKey"
	)
	suite.Run("Should return input found with valid bucket and key", func() {
		mockedHeadObject := &s3.HeadObjectInput{
			Bucket: &validBucket,
			Key:    &validKey,
		}
		suite.awsS3Mock.On("HeadObject", mockedHeadObject).Return(&s3.HeadObjectOutput{}, nil).Once()
		exists, err := suite.s3Client.CheckResourceExistance(validBucket, validKey)
		suite.Require().NoError(err, "Invoking method should not produce an internal error")
		suite.Require().Equal(true, exists, "Source input should return valid existance")
	})
	suite.Run("Should return input not found with invalid key", func() {
		mockedHeadObject := &s3.HeadObjectInput{
			Bucket: &validBucket,
			Key:    &invalidKey,
		}
		suite.awsS3Mock.On("HeadObject", mockedHeadObject).Return(&s3.HeadObjectOutput{}, mockedNotFoundError).Once()
		exists, err := suite.s3Client.CheckResourceExistance(validBucket, invalidKey)
		suite.Require().NoError(err, "Invoking method should not produce an internal error")
		suite.Require().Equal(false, exists, "Source input should return that it does not exist")
	})
	suite.Run("Should return input not found with invalid bucket", func() {
		mockedHeadObject := &s3.HeadObjectInput{
			Bucket: &validBucket,
			Key:    &invalidKey,
		}
		suite.awsS3Mock.On("HeadObject", mockedHeadObject).Return(&s3.HeadObjectOutput{}, mockedNotFoundError).Once()
		exists, err := suite.s3Client.CheckResourceExistance(validBucket, invalidKey)
		suite.Require().NoError(err, "Invoking method should not produce an internal error")
		suite.Require().Equal(false, exists, "Source input should return that it does not exist")
	})
	suite.Run("Should return an internal error with invalid region", func() {
		mockedHeadObject := &s3.HeadObjectInput{
			Bucket: &invalidBucket,
			Key:    &invalidKey,
		}
		suite.awsS3Mock.On("HeadObject", mockedHeadObject).Return(&s3.HeadObjectOutput{}, mockedUnexpectedError).Once()
		exists, err := suite.s3Client.CheckResourceExistance(invalidBucket, invalidKey)
		suite.Require().Equal(false, exists, "Source input should be empty if an internal error occurs")
		suite.Require().EqualError(errors.Cause(err), mockedUnexpectedError.Error(), "Invoking method should return expected internal error")
	})
}
