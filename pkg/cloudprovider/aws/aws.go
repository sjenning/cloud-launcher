package aws

import (
	"fmt"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/ec2"
	"github.com/sjenning/cloud-launcher/pkg/cloudprovider"
)

type Config struct {
	Region       string
	ImageID      string
	InstanceType string
	SubnetID     string
	KeyName      string
}

type Credentials struct {
	AccessKeyID     string
	SecretAccessKey string
}

type awsCloudProvider struct {
	Config
	svc         *ec2.EC2
	credentials Credentials
}

var _ cloudprovider.Interface = &awsCloudProvider{}

func NewCloudProvider(config *Config) (*awsCloudProvider, error) {
	session := session.New(&aws.Config{Region: aws.String(config.Region)})
	value, err := session.Config.Credentials.Get()
	if err != nil {
		return nil, err
	}
	svc := ec2.New(session)
	return &awsCloudProvider{
		Config: *config,
		svc:    svc,
		credentials: Credentials{
			AccessKeyID:     value.AccessKeyID,
			SecretAccessKey: value.SecretAccessKey,
		},
	}, nil

}

func (p *awsCloudProvider) CreateInstance() (string, error) {
	reservation, err := p.svc.RunInstances(&ec2.RunInstancesInput{
		ImageId:      aws.String(p.ImageID),
		InstanceType: aws.String(p.InstanceType),
		MinCount:     aws.Int64(1),
		MaxCount:     aws.Int64(1),
		SubnetId:     aws.String(p.SubnetID),
		KeyName:      aws.String(p.KeyName),
	})
	if err != nil {
		return "", err
	}
	return *(reservation.Instances[0].InstanceId), nil
}

func (p *awsCloudProvider) DeleteInstance(id string) error {
	input := &ec2.TerminateInstancesInput{
		InstanceIds: toAWSInstanceIDs(id),
	}
	_, err := p.svc.TerminateInstances(input)
	if err != nil {
		return err
	}
	return nil
}

func (p *awsCloudProvider) WaitForInstance(id string) {
	input := &ec2.DescribeInstancesInput{
		InstanceIds: toAWSInstanceIDs(id),
	}
	p.svc.WaitUntilInstanceRunning(input)
}

func (p *awsCloudProvider) TagInstance(id, key, value string) error {
	_, err := p.svc.CreateTags(&ec2.CreateTagsInput{
		Resources: toAWSInstanceIDs(id),
		Tags: []*ec2.Tag{
			{
				Key:   &key,
				Value: &value,
			},
		},
	})
	if err != nil {
		return err
	}
	return nil
}

func (p *awsCloudProvider) GetInstanceIP(id string) (string, error) {
	input := &ec2.DescribeInstancesInput{
		InstanceIds: toAWSInstanceIDs(id),
	}
	result, err := p.svc.DescribeInstances(input)
	if err != nil {
		return "", err
	}

	if len(result.Reservations) != 1 {
		return "", fmt.Errorf("reservation count in result was not exactly 1")
	}

	reservation := result.Reservations[0]
	if len(reservation.Instances) != 1 {
		return "", fmt.Errorf("instance count in result was not exactly 1")
	}

	instance := reservation.Instances[0]
	return *instance.PublicIpAddress, nil
}

func (p *awsCloudProvider) GetCredentials() interface{} {
	return p.credentials
}

func (p *awsCloudProvider) GetInstanceIDsByClusterName(name string) ([]string, error) {
	ids := make([]string, 0)
	params := &ec2.DescribeInstancesInput{
		Filters: []*ec2.Filter{
			&ec2.Filter{
				Name:   aws.String(fmt.Sprintf("tag:kubernetes.io/cluster/%s", name)),
				Values: []*string{aws.String("true")},
			},
		},
	}
	result, err := p.svc.DescribeInstances(params)
	if err != nil {
		return ids, err
	}
	for _, reservation := range result.Reservations {
		for _, instance := range reservation.Instances {
			ids = append(ids, *instance.InstanceId)
		}
	}
	return ids, nil
}

func toAWSInstanceIDs(id string) []*string {
	return []*string{&id}
}
