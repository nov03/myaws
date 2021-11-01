package myaws

import (
	"fmt"
	"regexp"
	"strings"

	"github.com/aws/aws-sdk-go/service/ec2"
)

// EC2SLsOptions customize the behavior of the Ls command.
type EC2SLsOptions struct {
	All       bool
	Quiet     bool
	FilterTag string
	Fields    []string
	Domain    []string
}

// EC2SLs describes EC2 snapshots.
func (client *Client) EC2SLs(options EC2SLsOptions) error {
	snapshots, err := client.FindEC2Snapshots(options.FilterTag, options.All)
	if err != nil {
		return err
	}

	for _, snapshot := range snapshots {
		fmt.Fprintln(client.stdout, formatEC2Snapshot(client, options, snapshot))
	}
	return nil
}

func formatEC2Snapshot(client *Client, options EC2SLsOptions, snapshot *ec2.Snapshot) string {
	formatFuncs := map[string]func(client *Client, options EC2SLsOptions, snapshot *ec2.Snapshot) string{
		"SnapshotId":        formatEC2SnapshotId,
		"Description":       formatEC2SnapshotDescription,
		"SnapshotStartTime": formatEC2SnapshotStartTime,
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
			value = formatEC2SnapshotTag(snapshot, key)

			if field == "Tag:Name" {
				address_name = value
			}

		} else {
			value = formatFuncs[field](client, options, snapshot)
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

func formatEC2SnapshotStartTime(client *Client, options EC2SLsOptions, snapshot *ec2.Snapshot) string {
	if snapshot.StartTime == nil {
		return "-"
	}
	return fmt.Sprintf("%-11s", *snapshot.StartTime)
}

func formatEC2SnapshotId(client *Client, options EC2SLsOptions, snapshot *ec2.Snapshot) string {
	return fmt.Sprintf("%-11s", *snapshot.SnapshotId)
}

func formatEC2SnapshotDescription(client *Client, options EC2SLsOptions, snapshot *ec2.Snapshot) string {
	return fmt.Sprintf("%-11s", *snapshot.Description)
}

func formatEC2SnapshotTag(snapshot *ec2.Snapshot, key string) string {
	var value string
	for _, t := range snapshot.Tags {
		if *t.Key == key {
			value = *t.Value
			break
		}
	}
	return value
}
