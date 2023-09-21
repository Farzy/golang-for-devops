package main

import (
	"context"
	"fmt"
	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/service/s3"
	"github.com/aws/aws-sdk-go-v2/service/s3/types"
	"os"
)

const region = "eu-west-3"
const bucketName = "farzad-aws-demo-test-bucket-346744"

func main() {
	var (
		s3Client *s3.Client
		err      error
	)

	ctx := context.Background()
	if s3Client, err = initS3Client(ctx, region); err != nil {
		fmt.Printf("initS3Client error: %s\n", err)
		os.Exit(1)
	}

	if err = createS3Bucket(ctx, s3Client, region); err != nil {
		fmt.Printf("createS3Bucket error: %s\n", err)
		os.Exit(1)
	}

	fmt.Printf("Bucket '%s' created\n", bucketName)
}

func initS3Client(ctx context.Context, region string) (*s3.Client, error) {
	cfg, err := config.LoadDefaultConfig(ctx, config.WithRegion(region))
	if err != nil {
		return nil, fmt.Errorf("unable to load SDK config, %v", err)
	}

	return s3.NewFromConfig(cfg), nil
}

func createS3Bucket(ctx context.Context, s3Client *s3.Client, region string) error {
	allBuckets, err := s3Client.ListBuckets(ctx, &s3.ListBucketsInput{})
	if err != nil {
		return fmt.Errorf("ListBuckers error: %s", err)
	}
	found := false
	for _, bucket := range allBuckets.Buckets {
		if *bucket.Name == bucketName {
			found = true
			break
		}
	}
	if !found {
		_, err = s3Client.CreateBucket(ctx, &s3.CreateBucketInput{
			Bucket: aws.String(bucketName),
			CreateBucketConfiguration: &types.CreateBucketConfiguration{
				LocationConstraint: types.BucketLocationConstraint(region),
			},
		})
		if err != nil {
			return fmt.Errorf("unable to create bucket: %s", err)
		}
	}
	return nil
}
