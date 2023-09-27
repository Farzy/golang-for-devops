package main

import (
	"context"
	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/s3"
	"github.com/aws/aws-sdk-go-v2/service/s3/types"
	"testing"
	"time"
)

type MockS3Client struct {
	ListBucketOutput   *s3.ListBucketsOutput
	CreateBucketOutput *s3.CreateBucketOutput
}

func (m *MockS3Client) ListBuckets(ctx context.Context, params *s3.ListBucketsInput, optFns ...func(*s3.Options)) (*s3.ListBucketsOutput, error) {
	return m.ListBucketOutput, nil
}

func (m *MockS3Client) CreateBucket(ctx context.Context, params *s3.CreateBucketInput, optFns ...func(*s3.Options)) (*s3.CreateBucketOutput, error) {
	return m.CreateBucketOutput, nil
}

func TestCreateS3Bucket_Creation(t *testing.T) {
	ctx := context.Background()
	err := createS3Bucket(ctx, &MockS3Client{
		ListBucketOutput: &s3.ListBucketsOutput{
			Buckets: []types.Bucket{
				{
					CreationDate: aws.Time(time.Now()),
					Name:         aws.String("test-bucket"),
				},
				{
					CreationDate: aws.Time(time.Now()),
					Name:         aws.String("test-bucket-2"),
				},
			},
		},
		CreateBucketOutput: &s3.CreateBucketOutput{
			Location: nil,
		},
	}, region)
	if err != nil {
		t.Fatalf("createS3Bucket error: %s", err)
	}
}

func TestCreateS3Bucket_Existing(t *testing.T) {
	ctx := context.Background()
	err := createS3Bucket(ctx, &MockS3Client{
		ListBucketOutput: &s3.ListBucketsOutput{
			Buckets: []types.Bucket{
				{
					CreationDate: aws.Time(time.Now()),
					Name:         aws.String("test-bucket"),
				},
				{
					CreationDate: aws.Time(time.Now()),
					Name:         aws.String(bucketName),
				},
			},
		},
		CreateBucketOutput: &s3.CreateBucketOutput{
			Location: nil,
		},
	}, region)
	if err != nil {
		t.Fatalf("createS3Bucket error: %s", err)
	}
}
