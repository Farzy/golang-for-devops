// Testing inspired by https://www.myhatchpad.com/insight/mocking-techniques-for-go/

package main

import (
	"context"
	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/feature/s3/manager"
	"github.com/aws/aws-sdk-go-v2/service/s3"
	"github.com/aws/aws-sdk-go-v2/service/s3/types"
	"testing"
	"time"
)

type MockS3Client struct {
	name                       string
	listBucketOutput           *s3.ListBucketsOutput
	createBucketOutput         *s3.CreateBucketOutput
	isCreateBucketCalled       bool
	expectedCreateBucketCalled bool
}

func (m *MockS3Client) ListBuckets(_ context.Context, _ *s3.ListBucketsInput, _ ...func(*s3.Options)) (*s3.ListBucketsOutput, error) {
	return m.listBucketOutput, nil
}

func (m *MockS3Client) CreateBucket(_ context.Context, _ *s3.CreateBucketInput, _ ...func(*s3.Options)) (*s3.CreateBucketOutput, error) {
	m.isCreateBucketCalled = true
	return m.createBucketOutput, nil
}

type MockS3Uploader struct {
	name string
}

func (m *MockS3Uploader) Upload(ctx context.Context, input *s3.PutObjectInput, opts ...func(*manager.Uploader)) (*manager.UploadOutput, error) {
	return &manager.UploadOutput{}, nil
}

func TestCreateS3Bucket_Creation(t *testing.T) {
	testTable := []MockS3Client{
		{
			name: "Bucket does not exist yet",
			listBucketOutput: &s3.ListBucketsOutput{
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
			createBucketOutput: &s3.CreateBucketOutput{
				Location: nil,
			},
			expectedCreateBucketCalled: true,
		},
		{
			name: "Bucket already exists",
			listBucketOutput: &s3.ListBucketsOutput{
				Buckets: []types.Bucket{
					{
						CreationDate: aws.Time(time.Now()),
						Name:         aws.String(bucketName),
					},
					{
						CreationDate: aws.Time(time.Now()),
						Name:         aws.String("test-bucket-2"),
					},
				},
			},
			createBucketOutput: &s3.CreateBucketOutput{
				Location: nil,
			},
			expectedCreateBucketCalled: false,
		},
	}
	ctx := context.Background()
	for _, tc := range testTable {
		t.Run(tc.name, func(t *testing.T) {
			err := createS3Bucket(ctx, &tc, region)
			if err != nil {
				t.Fatalf("createS3Bucket error: %s", err)
			}
			if tc.expectedCreateBucketCalled != tc.isCreateBucketCalled {
				t.Fatalf("expected isCreateBucketCalled to be %v, got %v",
					tc.expectedCreateBucketCalled,
					tc.isCreateBucketCalled)
			}
		})
	}
}

func TestUploadToS3Bucket(t *testing.T) {
	testTable := []MockS3Uploader{
		{
			name: "Upload file",
		},
	}
	ctx := context.Background()
	for _, tc := range testTable {
		t.Run(tc.name, func(t *testing.T) {
			err := uploadToS3Bucket(ctx, &tc)
			if err != nil {
				t.Fatalf("uploadToS3Bucket error: %s", err)
			}
		})
	}
}
