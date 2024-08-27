package storage

import (
	"bytes"
	"io"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/credentials"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/s3"
)

type S3Storage struct {
	session     *s3.S3
	credentials *credentials.Credentials
	bucket      string
}

type NewS3StorageInput struct {
	Endpoint       string
	Region         string
	Bucket         string
	DisableSSL     bool
	ForcePathStyle bool
	Credentials    *credentials.Credentials
}

func NewS3Storage(input NewS3StorageInput) *S3Storage {
	sess := session.Must(session.NewSession(&aws.Config{
		Endpoint:         aws.String(input.Endpoint),
		Region:           aws.String(input.Region),
		DisableSSL:       aws.Bool(input.DisableSSL),
		S3ForcePathStyle: aws.Bool(input.ForcePathStyle),
		Credentials:      input.Credentials,
	}))
	return &S3Storage{
		session:     s3.New(sess),
		credentials: input.Credentials,
		bucket:      input.Bucket,
	}
}

func (s *S3Storage) Store(key string, data []byte, metadata map[string]*string) error {
	_, err := s.session.PutObject(&s3.PutObjectInput{
		Bucket:   aws.String(s.bucket),
		Key:      aws.String(key),
		Body:     bytes.NewReader(data),
		Metadata: metadata,
	})
	return err
}

func (s *S3Storage) Retrieve(key string) ([]byte, error) {
	result, err := s.session.GetObject(&s3.GetObjectInput{
		Bucket: aws.String(s.bucket),
		Key:    aws.String(key),
	})
	if err != nil {
		return nil, err
	}
	defer result.Body.Close()
	return io.ReadAll(result.Body)
}

func (s *S3Storage) Delete(key string) error {
	_, err := s.session.DeleteObject(&s3.DeleteObjectInput{
		Bucket: aws.String(s.bucket),
		Key:    aws.String(key),
	})
	return err
}
