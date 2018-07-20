package myaws

import (
	"time"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/request"
	"github.com/aws/aws-sdk-go/service/autoscaling"
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
	return client.waitUntilGroupAllInstancesInServiceWithContext(aws.BackgroundContext(), input)
}

// waitUntilAutoScalingGroupAllInstancesInService waits until all instances are
// in service with context.
// Since the `(autoscaling.*AutoScaling) WaitUntilGroupInServiceWithContext`
// checks `>=MinSize` and it is not suitable for detaching. So we implement a
// customized waiter here.  When the number of desired instances increase or
// decrease, the affected instances are in states other than InService until
// the operation completes. So we can simply check that all the states of
// instances are InService.
func (client *Client) waitUntilGroupAllInstancesInServiceWithContext(ctx aws.Context, input *autoscaling.DescribeAutoScalingGroupsInput, opts ...request.WaiterOption) error {
	w := request.Waiter{
		Name:        "WaitUntilGroupAllInstancesInService",
		MaxAttempts: 40,
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
