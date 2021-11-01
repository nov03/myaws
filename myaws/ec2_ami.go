package myaws

import (
	"fmt"
	"regexp"
	"strings"

	"github.com/aws/aws-sdk-go/service/ec2"
)

// EC2ALsOptions customize the behavior of the Ls command.
type EC2ALsOptions struct {
	All       bool
	Quiet     bool
	FilterTag string
	Fields    []string
	Domain    []string
}

// EC2ALs describes EC2 images.
func (client *Client) EC2ALs(options EC2ALsOptions) error {
	images, err := client.FindEC2Amis(options.FilterTag, options.All)
	if err != nil {
		return err
	}

	for _, image := range images {
		fmt.Fprintln(client.stdout, formatEC2Ami(client, options, image))
	}
	return nil
}

func formatEC2Ami(client *Client, options EC2ALsOptions, image *ec2.Image) string {
	formatFuncs := map[string]func(client *Client, options EC2ALsOptions, image *ec2.Image) string{
		"AmiName":      formatEC2AmiName,
		"ImageId":      formatEC2AmiImageId,
		"CreationDate": formatEC2AmiCreationDate,
	}

	var outputFields []string
	if options.Quiet {
		outputFields = []string{"AllocationId"}
	} else {
		outputFields = options.Fields
	}

	output := []string{}
	address_name := ""

	for _, field := range outputFields {
		value := ""

		if strings.Index(field, "Tag:") != -1 {
			key := strings.Split(field, ":")[1]
			value = formatEC2AMITag(image, key)

			if field == "Tag:Name" {
				address_name = value
			}

		} else {
			value = formatFuncs[field](client, options, image)
		}
		output = append(output, value)
	}

	// -Dの指定時に、引数で指定された文字列とインスタンスのtag:nameに合致する対象だけoutputに追記する。
	if len(options.Domain) == 0 {
		return strings.Join(output[:], "\t")
	} else {
		if regexp.MustCompile(options.Domain[0]).Match([]byte(address_name)) {
			return strings.Join(output[:], "\t")
		} else {
			output = nil
			return strings.Join(output[:], "")
		}
	}
}

func formatEC2AmiName(client *Client, options EC2ALsOptions, image *ec2.Image) string {
	if image.Name == nil {
		return "-"
	}
	return fmt.Sprintf("%-11s", *image.Name)
}

func formatEC2AmiImageId(client *Client, options EC2ALsOptions, image *ec2.Image) string {
	return fmt.Sprintf("%-11s", *image.ImageId)
}

func formatEC2AmiCreationDate(client *Client, options EC2ALsOptions, image *ec2.Image) string {
	return fmt.Sprintf("%-11s", *image.CreationDate)
}

func formatEC2AMITag(image *ec2.Image, key string) string {
	var value string
	for _, t := range image.Tags {
		if *t.Key == key {
			value = *t.Value
			break
		}
	}
	return value
}
