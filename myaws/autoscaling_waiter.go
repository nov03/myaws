package myaws

import (
	"fmt"
	"time"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/request"
	"github.com/aws/aws-sdk-go/service/autoscaling"
	"github.com/k0kubun/pp"
	"github.com/pkg/errors"
)

// waitUntilAutoScalingGroupDesiredState is a helper function which waits until
// the AutoScaling Group converges to the desired state.  We only check the
// status of AutoScaling Group.  If the ASG has an ELB, the health check status
// of ELB can link with the health status of ASG, so we don't check the status
// of ELB here.
func (client *Client) waitUntilAutoScalingGroupDesiredState(asgName string) error {
	err := client.waitUntilAutoScalingGroupAllInstancesInService(
		&autoscaling.DescribeAutoScalingGroupsInput{
			AutoScalingGroupNames: []*string{&asgName},
		},
	)
	if err != nil {
		return errors.Wrap(err, "waitUntilAutoScalingGroupAllInstancesInService failed:")
	}

	return nil
}

// waitUntilAutoScalingGroupAllInstancesInService waits until all instances are
// in service.
func (client *Client) waitUntilAutoScalingGroupAllInstancesInService(input *autoscaling.DescribeAutoScalingGroupsInput) error {
	ctx := aws.BackgroundContext()

	if err := client.waitUntilAutoScalingGroupNumberOfInstancesIsDesiredCapacityWithContext(ctx, input); err != nil {
		return err
	}

	return client.waitUntilAutoScalingGroupAllInstancesInServiceWithContext(ctx, input)
}

func (client *Client) waitUntilAutoScalingGroupNumberOfInstancesIsDesiredCapacityWithContext(ctx aws.Context, input *autoscaling.DescribeAutoScalingGroupsInput, opts ...request.WaiterOption) error {
	asgName := *(input.AutoScalingGroupNames[0])
	desiredCapacity, err := client.getAutoScalingGroupDesiredCapacity(asgName)
	if err != nil {
		return err
	}

	matcher := fmt.Sprintf("AutoScalingGroups[].[length(Instances) == `%d`][]", desiredCapacity)
	pp.Printf("matcher: %+v\n", matcher)

	w := request.Waiter{
		Name:        "WaitUntilAutoScalingGroupNumberOfInstancesIsDesiredCapacity",
		MaxAttempts: 20,
		Delay:       request.ConstantWaiterDelay(15 * time.Second),
		Acceptors: []request.WaiterAcceptor{
			{
				State:   request.SuccessWaiterState,
				Matcher: request.PathAllWaiterMatch, Argument: matcher,
				Expected: true,
			},
		},
		Logger: client.config.Logger,
		NewRequest: func(opts []request.Option) (*request.Request, error) {
			var inCpy *autoscaling.DescribeAutoScalingGroupsInput
			if input != nil {
				tmp := *input
				inCpy = &tmp
			}
			req, _ := client.AutoScaling.DescribeAutoScalingGroupsRequest(inCpy)
			req.SetContext(ctx)
			req.ApplyOptions(opts...)
			return req, nil
		},
	}
	w.ApplyOptions(opts...)

	return w.WaitWithContext(ctx)
}

func (client *Client) getAutoScalingGroupDesiredCapacity(asgName string) (int64, error) {
	input := &autoscaling.DescribeAutoScalingGroupsInput{
		AutoScalingGroupNames: []*string{&asgName},
	}

	response, err := client.AutoScaling.DescribeAutoScalingGroups(input)
	if err != nil {
		return 0, errors.Wrap(err, "getAutoScalingGroupDesiredCapacity failed:")
	}

	desiredCapacity := response.AutoScalingGroups[0].DesiredCapacity

	pp.Printf("desiredCapacity: %+v\n", desiredCapacity)
	return *desiredCapacity, nil
}

// waitUntilAutoScalingGroupAllInstancesInService waits until all instances are
// in service with context.
// Since the `(autoscaling.*AutoScaling) WaitUntilGroupInServiceWithContext`
// checks `>=MinSize` and it is not suitable for detaching. So we implement a
// customized waiter here.  When the number of desired instances increase or
// decrease, the affected instances are in states other than InService until
// the operation completes. So we can simply check that all the states of
// instances are InService.
func (client *Client) waitUntilAutoScalingGroupAllInstancesInServiceWithContext(ctx aws.Context, input *autoscaling.DescribeAutoScalingGroupsInput, opts ...request.WaiterOption) error {
	w := request.Waiter{
		Name:        "WaitUntilAutoScalingGroupAllInstancesInService",
		MaxAttempts: 20,
		Delay:       request.ConstantWaiterDelay(15 * time.Second),
		Acceptors: []request.WaiterAcceptor{
			{
				State:   request.SuccessWaiterState,
				Matcher: request.PathAllWaiterMatch, Argument: "AutoScalingGroups[].Instances[].LifecycleState",
				Expected: "InService",
			},
		},
		Logger: client.config.Logger,
		NewRequest: func(opts []request.Option) (*request.Request, error) {
			var inCpy *autoscaling.DescribeAutoScalingGroupsInput
			if input != nil {
				tmp := *input
				inCpy = &tmp
			}
			req, _ := client.AutoScaling.DescribeAutoScalingGroupsRequest(inCpy)
			req.SetContext(ctx)
			req.ApplyOptions(opts...)
			return req, nil
		},
	}
	w.ApplyOptions(opts...)

	return w.WaitWithContext(ctx)
}
