package main

import (
	"context"
	"fmt"
	"github.com/aws/aws-sdk-go-v2/service/ec2/types"
	"os"
	"strings"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/service/ec2"
)

const (
	privateSSHKeyFile = "go-aws-ec2.pem"
	keyPairName       = "go-aws-demo"
)

func main() {
	var (
		instanceId string
		err        error
	)
	ctx := context.Background()

	if instanceId, err = createEC2(ctx, "eu-west-3"); err != nil {
		fmt.Printf("createEC2 errors: %s", err)
		os.Exit(1)
	}

	fmt.Printf("Instance id: %s\n", instanceId)
}

func createEC2(ctx context.Context, region string) (string, error) {
	// Using the SDK's default configuration, loading additional config
	// and credentials values from the environment variables, shared
	// credentials, and shared configuration files
	cfg, err := config.LoadDefaultConfig(ctx, config.WithRegion(region))
	if err != nil {
		return "", fmt.Errorf("unable to load SDK config, %v", err)
	}

	ec2Client := ec2.NewFromConfig(cfg)

	pairs, err := ec2Client.DescribeKeyPairs(ctx, &ec2.DescribeKeyPairsInput{
		KeyNames: []string{keyPairName},
	})
	if err != nil && !strings.Contains(err.Error(), "InvalidKeyPair.NotFound") {
		return "", fmt.Errorf("DescribeKeyPairs error: %s", err)
	}
	if pairs == nil || len(pairs.KeyPairs) == 0 {
		keyPairOutput, err := ec2Client.CreateKeyPair(ctx, &ec2.CreateKeyPairInput{
			KeyName: aws.String(keyPairName),
		})
		if err != nil {
			return "", fmt.Errorf("CreateKeyPair error: %s", err)
		}
		fmt.Printf("Writing private SSH key to file %s\n", privateSSHKeyFile)
		err = os.WriteFile(privateSSHKeyFile, []byte(*keyPairOutput.KeyMaterial), 0600)
		if err != nil {
			return "", fmt.Errorf("WriteFile: cannot with private ssh key to file: %s\n", err)
		}
	}
	imagesOutput, err := ec2Client.DescribeImages(ctx, &ec2.DescribeImagesInput{
		Filters: []types.Filter{
			{
				Name:   aws.String("name"),
				Values: []string{"ubuntu/images/hvm-ssd/ubuntu-focal-20.04-amd64-server-*"},
			},
			{
				Name:   aws.String("virtualization-type"),
				Values: []string{"hvm"},
			},
		},
		Owners: []string{"099720109477"},
	})
	if err != nil {
		return "", fmt.Errorf("DescribeImage error: %s", err)
	}

	if len(imagesOutput.Images) == 0 {
		return "", fmt.Errorf("imageOutput.Images is of length 0")
	}

	instances, err := ec2Client.RunInstances(
		ctx,
		&ec2.RunInstancesInput{
			MaxCount:     aws.Int32(1),
			MinCount:     aws.Int32(1),
			ImageId:      imagesOutput.Images[0].ImageId,
			InstanceType: types.InstanceTypeT3Micro,
			KeyName:      aws.String("go-aws-demo"),
			SubnetId:     aws.String("subnet-0d458afb2d2c222f6"),
		},
	)
	if err != nil {
		return "", fmt.Errorf("RunInstances error: %s", err)
	}

	if len(instances.Instances) == 0 {
		return "", fmt.Errorf("instance.Instances is of length 0")
	}

	return *instances.Instances[0].InstanceId, nil
}
