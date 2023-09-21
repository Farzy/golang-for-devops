package main

import (
	"context"
	"fmt"
	"os"
	"strings"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/feature/s3/manager"
	"github.com/aws/aws-sdk-go-v2/service/s3"
	"github.com/aws/aws-sdk-go-v2/service/s3/types"
)

const region = "eu-west-3"
const bucketName = "farzad-aws-demo-test-bucket-346744"
const poemFileName = "ode-to-go.txt"

func main() {
	var (
		s3Client *s3.Client
		err      error
		out      []byte
		outFile  string
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

	if err = uploadToS3Bucket(ctx, s3Client); err != nil {
		fmt.Printf("uploadToS3Bucket error: %s\n", err)
		os.Exit(1)
	}
	fmt.Printf("Upload complete!\n")

	if out, err = readFromS3Bucket(ctx, s3Client); err != nil {
		fmt.Printf("readFromS3Bucket error: %s\n", err)
		os.Exit(1)
	}
	fmt.Printf("Read complete: %s\n", out)

	if outFile, err = downloadFromS3Bucket(ctx, s3Client); err != nil {
		fmt.Printf("downloadFromS3Bucket error: %s\n", err)
		os.Exit(1)
	}
	if out, err = os.ReadFile(outFile); err != nil {
		fmt.Printf("ReadFile error: %s\n", err)
		os.Exit(1)
	}
	fmt.Printf("Download complete: %s\n", out)
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

func uploadToS3Bucket(ctx context.Context, s3Client *s3.Client) error {
	uploader := manager.NewUploader(s3Client)
	_, err := uploader.Upload(ctx, &s3.PutObjectInput{
		Bucket: aws.String(bucketName),
		Key:    aws.String("directory/test.txt"),
		Body:   strings.NewReader("Hello world!"),
	})
	if err != nil {
		return fmt.Errorf("upload error:, %v", err)
	}

	file, err := os.Open(poemFileName)
	if err != nil {
		return fmt.Errorf("cannot open file '%': %s", poemFileName, err)
	}
	_, err = uploader.Upload(ctx, &s3.PutObjectInput{
		Bucket: aws.String(bucketName),
		Key:    aws.String("directory2/ode-to-go.txt"),
		Body:   file,
	})
	if err != nil {
		return fmt.Errorf("upload error:, %v", err)
	}

	return nil
}

func readFromS3Bucket(ctx context.Context, s3Client *s3.Client) ([]byte, error) {
	buffer := manager.NewWriteAtBuffer([]byte{})
	downloader := manager.NewDownloader(s3Client)
	numBytes, err := downloader.Download(ctx, buffer, &s3.GetObjectInput{
		Bucket: aws.String(bucketName),
		Key:    aws.String("directory2/ode-to-go.txt"),
	})
	if err != nil {
		return nil, fmt.Errorf("download error:, %v", err)
	}

	if numBytesReceived := len(buffer.Bytes()); numBytes != int64(numBytesReceived) {
		return nil, fmt.Errorf("numbytes received doesn't match: %d vs %d", numBytes, numBytesReceived)
	}

	return buffer.Bytes(), nil
}

func downloadFromS3Bucket(ctx context.Context, s3Client *s3.Client) (string, error) {
	buffer := manager.NewWriteAtBuffer([]byte{})
	downloader := manager.NewDownloader(s3Client)
	numBytes, err := downloader.Download(ctx, buffer, &s3.GetObjectInput{
		Bucket: aws.String(bucketName),
		Key:    aws.String("directory/test.txt"),
	})
	if err != nil {
		return "", fmt.Errorf("download error:, %v", err)
	}

	if numBytesReceived := len(buffer.Bytes()); numBytes != int64(numBytesReceived) {
		return "", fmt.Errorf("numbytes received doesn't match: %d vs %d", numBytes, numBytesReceived)
	}

	err = os.WriteFile("test-dl.txt", buffer.Bytes(), 0600)
	if err != nil {
		return "", fmt.Errorf("WriteFile error: %s", err)
	}

	return "test-dl.txt", nil
}
