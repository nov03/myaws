package myaws

import (
	"fmt"
	"regexp"
	"strings"

	"github.com/aws/aws-sdk-go/service/ec2"
)

// EC2ILsOptions customize the behavior of the Ls command.
type EC2ILsOptions struct {
	All       bool
	Quiet     bool
	FilterTag string
	Fields    []string
	Domain    []string
}

// EC2ILs describes EC2 addresses.
func (client *Client) EC2ILs(options EC2ILsOptions) error {
	addresses, err := client.FindEC2Ips(options.FilterTag, options.All)
	if err != nil {
		return err
	}

	for _, address := range addresses {
		fmt.Fprintln(client.stdout, formatEC2Ip(client, options, address))
	}
	return nil
}

func formatEC2Ip(client *Client, options EC2ILsOptions, address *ec2.Address) string {
	formatFuncs := map[string]func(client *Client, options EC2ILsOptions, address *ec2.Address) string{
		"PublicIp":         formatEC2PublicIp,
		"AllocationId":     formatEC2AllocationId,
		"InstanceId":       formatEC2InstanceId,
		"PrivateIpAddress": formatEC2PrivateIpAddress,
		"AssociationId":    formatEC2AssociationId,
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
			value = formatEC2ITag(address, key)

			if field == "Tag:Name" {
				address_name = value
			}

		} else {
			value = formatFuncs[field](client, options, address)
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

func formatEC2PublicIp(client *Client, options EC2ILsOptions, address *ec2.Address) string {
	return fmt.Sprintf("%-11s", *address.PublicIp)
}

func formatEC2AllocationId(client *Client, options EC2ILsOptions, address *ec2.Address) string {
	return fmt.Sprintf("%-11s", *address.AllocationId)
}

func formatEC2InstanceId(client *Client, options EC2ILsOptions, address *ec2.Address) string {
	if address.InstanceId == nil {
		return "-"
	}
	return fmt.Sprintf("%-11s", *address.InstanceId)
}

func formatEC2PrivateIpAddress(client *Client, options EC2ILsOptions, address *ec2.Address) string {
	return fmt.Sprintf("%-11s", *address.PrivateIpAddress)
}

func formatEC2AssociationId(client *Client, options EC2ILsOptions, address *ec2.Address) string {
	return fmt.Sprintf("%-11s", *address.AssociationId)
}

func formatEC2ITag(address *ec2.Address, key string) string {
	var value string
	for _, t := range address.Tags {
		if *t.Key == key {
			value = *t.Value
			break
		}
	}
	return value
}
