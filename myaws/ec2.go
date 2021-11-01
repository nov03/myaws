package myaws

import (
	"strings"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/ec2"
	"github.com/pkg/errors"
)

// FindEC2Instances returns an array of instances matching the conditions.
func (client *Client) FindEC2Instances(filterTag string, all bool) ([]*ec2.Instance, error) {
	params := &ec2.DescribeInstancesInput{
		Filters: []*ec2.Filter{
			buildEC2StateFilter(all),
			buildEC2TagFilter(filterTag),
		},
	}

	response, err := client.EC2.DescribeInstances(params)
	if err != nil {
		return nil, errors.Wrap(err, "DescribeInstances failed")
	}

	var instances []*ec2.Instance
	for _, reservation := range response.Reservations {
		for _, instance := range reservation.Instances {
			instances = append(instances, instance)
		}
	}

	return instances, nil
}

func buildEC2StateFilter(all bool) *ec2.Filter {
	var stateFilter *ec2.Filter
	if !all {
		stateFilter = &ec2.Filter{
			Name: aws.String("instance-state-name"),
			Values: []*string{
				aws.String("running"),
			},
		}
	}
	return stateFilter
}

func buildEC2TagFilter(filterTag string) *ec2.Filter {
	var tagFilter *ec2.Filter
	if filterTag != "" {
		tagParts := strings.Split(filterTag, ":")
		tagFilter = &ec2.Filter{
			Name: aws.String("tag:" + tagParts[0]),
			Values: []*string{
				aws.String("*" + tagParts[1] + "*"),
			},
		}
	}
	return tagFilter
}

// FindEC2Volumes returns an array of volumes matching the conditions.
func (client *Client) FindEC2Volumes(filterTag string, all bool) ([]*ec2.Volume, error) {
	params := &ec2.DescribeVolumesInput{}

	response, err := client.EC2.DescribeVolumes(params)
	if err != nil {
		return nil, errors.Wrap(err, "DescribeVolumes failed")
	}

	var volumes []*ec2.Volume

	for _, volume := range response.Volumes {
		volumes = append(volumes, volume)
	}

	return volumes, nil
}

// FindEC2Ips returns an array of volumes matching the conditions.
func (client *Client) FindEC2Ips(filterTag string, all bool) ([]*ec2.Address, error) {
	params := &ec2.DescribeAddressesInput{}

	response, err := client.EC2.DescribeAddresses(params)
	if err != nil {
		return nil, errors.Wrap(err, "DescribeIps failed")
	}

	var Addresses []*ec2.Address

	for _, Address := range response.Addresses {
		Addresses = append(Addresses, Address)
	}

	return Addresses, nil
}

// FindEC2Amis returns an array of volumes matching the conditions.
func (client *Client) FindEC2Amis(filterTag string, all bool) ([]*ec2.Image, error) {
	params := &ec2.DescribeImagesInput{}

	response, err := client.EC2.DescribeImages(params)
	if err != nil {
		return nil, errors.Wrap(err, "DescribeAmis failed")
	}

	var Images []*ec2.Image

	for _, Image := range response.Images {
		Images = append(Images, Image)
	}

	return Images, nil
}
