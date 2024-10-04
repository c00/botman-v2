package storageprovider

import (
	"bytes"
	"context"
	"errors"
	"fmt"
	"io"
	"path"

	awsConfig "github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/credentials"
	"github.com/aws/aws-sdk-go-v2/service/s3"
	"github.com/aws/aws-sdk-go-v2/service/s3/types"
)

var _ StorageProvider = (*S3Storage)(nil)

func NewS3Store(conf S3Config) *S3Storage {
	return &S3Storage{
		Config: conf,
	}
}

type S3Storage struct {
	Config S3Config
	client *s3.Client
}

func (s *S3Storage) getClient() (*s3.Client, error) {
	if s.client != nil {
		return s.client, nil
	}

	awsCfg, err := awsConfig.LoadDefaultConfig(context.Background(),
		awsConfig.WithRegion(s.Config.Region),
	)

	if err != nil {
		return nil, err
	}

	s.client = s3.NewFromConfig(awsCfg, func(o *s3.Options) {
		o.BaseEndpoint = &s.Config.Endpoint
		o.Credentials = credentials.NewStaticCredentialsProvider(s.Config.AccessKey, s.Config.SecretKey, "")
	})

	return s.client, nil
}

func (s *S3Storage) getFullPath(name string) string {
	return path.Join(s.Config.BasePath, name)
}

func (s *S3Storage) Save(name string, data []byte) (string, error) {
	if name == "" {
		return "", ErrEmptyFilename
	}

	client, err := s.getClient()
	if err != nil {
		return "", fmt.Errorf("could not get s3 client: %w", err)
	}

	key := s.getFullPath(name)

	_, err = client.PutObject(context.Background(), &s3.PutObjectInput{
		Body:   bytes.NewReader(data),
		Bucket: &s.Config.Bucket,
		Key:    &key,
		ACL:    types.ObjectCannedACLPublicRead,
	})
	if err != nil {
		return "", fmt.Errorf("upload to s3 failed: %w", err)
	}

	//build public url.
	publicUrl := path.Join(s.Config.Endpoint, s.Config.Bucket, key)

	return publicUrl, nil
}

func (s *S3Storage) Load(name string) ([]byte, error) {
	if name == "" {
		return nil, ErrEmptyFilename
	}

	client, err := s.getClient()
	if err != nil {
		return nil, fmt.Errorf("could not get s3 client: %w", err)
	}

	key := s.getFullPath(name)

	result, err := client.GetObject(context.Background(), &s3.GetObjectInput{
		Bucket: &s.Config.Bucket,
		Key:    &key,
	})
	if err != nil {
		return nil, err
	}

	data, err := io.ReadAll(result.Body)
	if err != nil {
		return nil, fmt.Errorf("body read error: %w", err)
	}

	return data, nil
}

func (s *S3Storage) DeleteFile(name string) error {
	if name == "" {
		return ErrEmptyFilename
	}

	client, err := s.getClient()
	if err != nil {
		return fmt.Errorf("could not get s3 client: %w", err)
	}

	key := s.getFullPath(name)

	_, err = client.DeleteObject(context.Background(), &s3.DeleteObjectInput{
		Bucket: &s.Config.Bucket,
		Key:    &key,
	})

	//Note, no errors are thrown if the key already doesn't exist.
	if err != nil {
		return err
	}

	return nil
}

func (s *S3Storage) DeleteFolder(prefix string) error {
	client, err := s.getClient()
	if err != nil {
		return fmt.Errorf("could not get s3 client: %w", err)
	}

	folder := path.Join(s.Config.BasePath, prefix)

	list, err := client.ListObjectsV2(context.Background(), &s3.ListObjectsV2Input{
		Bucket: &s.Config.Bucket,
		Prefix: &folder,
	})
	if err != nil {
		return fmt.Errorf("could not list items in folder: %w", err)
	}

	objects := []types.ObjectIdentifier{}

	for _, c := range list.Contents {
		objects = append(objects, types.ObjectIdentifier{Key: c.Key})
	}

	_, err = client.DeleteObjects(context.Background(), &s3.DeleteObjectsInput{
		Bucket: &s.Config.Bucket,
		Delete: &types.Delete{Objects: objects},
	})

	if err != nil {
		return fmt.Errorf("could not delete folder: %w", err)
	}

	return nil
}

func (s *S3Storage) HasFile(name string) (bool, error) {
	client, err := s.getClient()
	if err != nil {
		return false, fmt.Errorf("could not get s3 client: %w", err)
	}

	key := s.getFullPath(name)

	_, err = client.HeadObject(context.Background(), &s3.HeadObjectInput{
		Bucket: &s.Config.Bucket,
		Key:    &key,
	})
	if err != nil {
		var notFound *types.NotFound
		if errors.As(err, &notFound) {
			return false, nil
		}
		return false, fmt.Errorf("could not head item: %w", err)
	}
	return true, nil
}

type S3Config struct {
	AccessKey      string `yaml:"accessKey"`
	SecretKey      string `yaml:"secretKey"`
	BasePath       string `yaml:"basePath"`
	Bucket         string `yaml:"bucket"`
	Endpoint       string `yaml:"endpoint"`
	ForcePathStyle bool   `yaml:"forcePathStyle"`
	Region         string `yaml:"region"`
}
