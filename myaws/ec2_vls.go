package myaws

import (
	"fmt"
	"regexp"
	"strings"

	"github.com/aws/aws-sdk-go/service/ec2"
)

// EC2VLsOptions customize the behavior of the Ls command.
type EC2VLsOptions struct {
	All       bool
	Quiet     bool
	FilterTag string
	Fields    []string
	Domain    []string
}

// EC2VLs describes EC2 volumes.
func (client *Client) EC2VLs(options EC2VLsOptions) error {
	volumes, err := client.FindEC2Volumes(options.FilterTag, options.All)
	if err != nil {
		return err
	}

	for _, volume := range volumes {
		fmt.Fprintln(client.stdout, formatEC2Volume(client, options, volume))
	}
	return nil
}

func formatEC2Volume(client *Client, options EC2VLsOptions, volume *ec2.Volume) string {
	formatFuncs := map[string]func(client *Client, options EC2VLsOptions, volume *ec2.Volume) string{
		"VolumeId":    formatEC2VolumeID,
		"VolumeType":  formatEC2VolumeType,
		"Size":        formatEC2VolumeSize,
		"Attachments": formatEC2VolumeAttachments,
	}

	var outputFields []string
	if options.Quiet {
		outputFields = []string{"VolumeId"}
	} else {
		outputFields = options.Fields
	}

	output := []string{}
	volume_name := ""

	for _, field := range outputFields {
		value := ""

		if strings.Index(field, "Tag:") != -1 {
			key := strings.Split(field, ":")[1]
			value = formatEC2VTag(volume, key)

			if field == "Tag:Name" {
				volume_name = value
			}

		} else {
			value = formatFuncs[field](client, options, volume)
		}
		output = append(output, value)
	}

	// -Dの指定時に、引数で指定された文字列とインスタンスのtag:nameに合致する対象だけoutputに追記する。
	if len(options.Domain) == 0 {
		return strings.Join(output[:], "\t")
	} else {
		if regexp.MustCompile(options.Domain[0]).Match([]byte(volume_name)) {
			return strings.Join(output[:], "\t")
		} else {
			output = nil
			return strings.Join(output[:], "")
		}
	}
}

func formatEC2VolumeID(client *Client, options EC2VLsOptions, volume *ec2.Volume) string {
	return *volume.VolumeId
}

func formatEC2VolumeType(client *Client, options EC2VLsOptions, volume *ec2.Volume) string {
	return fmt.Sprintf("%-11s", *volume.VolumeType)
}

func formatEC2VolumeSize(client *Client, options EC2VLsOptions, volume *ec2.Volume) string {
	return fmt.Sprintf("%dGib", *volume.Size)
}

func formatEC2VolumeAttachments(client *Client, options EC2VLsOptions, volume *ec2.Volume) string {

	return fmt.Sprintf("%-11s\t%-11s", *volume.Attachments[0].InstanceId, *volume.Attachments[0].Device)
}

func formatEC2VTag(volume *ec2.Volume, key string) string {
	var value string
	for _, t := range volume.Tags {
		if *t.Key == key {
			value = *t.Value
			break
		}
	}
	return value
}
