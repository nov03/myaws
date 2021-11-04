package myaws

import (
	"fmt"
	"strings"

	"github.com/aws/aws-sdk-go/service/elbv2"
	"github.com/pkg/errors"
)

// EC2LsOptions customize the behavior of the Ls command.
type ELBv2Options struct {
	All       bool
	Quiet     bool
	FilterTag string
	Fields    []string
	Domain    []string
}

// ELBV2Ls describes ELBV2s.
func (client *Client) ELBV2Ls(options ELBv2Options) error {
	params := &elbv2.DescribeLoadBalancersInput{}

	response, err := client.ELBV2.DescribeLoadBalancers(params)
	if err != nil {
		return errors.Wrap(err, "DescribeLoadBalancers failed:")
	}

	for _, lb := range response.LoadBalancers {
		fmt.Fprintln(client.stdout, formatLoadBalancerV2(lb))
	}

	return nil
}

func formatLoadBalancerV2(lb *elbv2.LoadBalancer) string {
	var a string
	n := len(lb.AvailabilityZones)

	output := []string{
		*lb.LoadBalancerName,
		*lb.DNSName,
		*lb.VpcId,
		*lb.Type,
	}
	// LBの振り分け先をoutputに格納、振り分け先のリージョンはVPCの右列に出力させたいがやり方がわからない
	for i := 0; i < n; i++ {
		a = *lb.AvailabilityZones[i].ZoneName
		output = append(output, a)
	}
	// a = len(options.Domain)
	// options.Domain[0]
	// if regexp.MustCompile(options.Domain[0]).Match([]byte(*lb.LoadBalancerName)) {
	// 	return strings.Join(output[:], "\t")
	// } else {
	// 	output = nil
	// 	return strings.Join(output[:], "")
	// }
	return strings.Join(output[:], "\t")
}
