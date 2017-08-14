// Copyright © 2017 Axel Springer SE
//
// Permission is hereby granted, free of charge, to any person obtaining a copy
// of this software and associated documentation files (the "Software"), to deal
// in the Software without restriction, including without limitation the rights
// to use, copy, modify, merge, publish, distribute, sublicense, and/or sell
// copies of the Software, and to permit persons to whom the Software is
// furnished to do so, subject to the following conditions:
//
// The above copyright notice and this permission notice shall be included in
// all copies or substantial portions of the Software.
//
// THE SOFTWARE IS PROVIDED "AS IS", WITHOUT WARRANTY OF ANY KIND, EXPRESS OR
// IMPLIED, INCLUDING BUT NOT LIMITED TO THE WARRANTIES OF MERCHANTABILITY,
// FITNESS FOR A PARTICULAR PURPOSE AND NONINFRINGEMENT. IN NO EVENT SHALL THE
// AUTHORS OR COPYRIGHT HOLDERS BE LIABLE FOR ANY CLAIM, DAMAGES OR OTHER
// LIABILITY, WHETHER IN AN ACTION OF CONTRACT, TORT OR OTHERWISE, ARISING FROM,
// OUT OF OR IN CONNECTION WITH THE SOFTWARE OR THE USE OR OTHER DEALINGS IN
// THE SOFTWARE.

package aws

import (
	amzn "github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/ec2metadata"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/autoscaling"
	"github.com/aws/aws-sdk-go/service/ec2"
)

// EC2 wraps the interfaces for the instance
type EC2 struct {
	*session.Session
	*ec2.EC2
	*ec2metadata.EC2Metadata
	*autoscaling.AutoScaling
}

// NewEC2 creates a new instance of EC2Instance to interface AWS
func NewEC2(session *session.Session) (*EC2, error) {
	// use the session first to determine the region for the new session with aws
	metadata := ec2metadata.New(session)
	identity, err := metadata.GetInstanceIdentityDocument()

	// error, if not instance identity is accessible
	if err != nil || identity.Region == "" {
		return nil, ErrEC2MetadataNotAvailable
	}

	// create new session with detected region
	session = NewSession(&amzn.Config{
		Region: &identity.Region,
	})

	// use the new session to create clients
	return &EC2{session, ec2.New(session), ec2metadata.New(session), autoscaling.New(session)}, nil
}

// GetAutoScalingGroupPeers returns the peers of the EC2 instance in the autoscaling group
func (e *EC2) GetAutoScalingGroupPeers() ([]*ec2.Instance, error) {
	// get current identity document
	instanceIdentity, err := e.GetInstanceIdentityDocument()

	// if err in client return a custom error in interface
	if err != nil {
		return nil, ErrEC2MetadataNotAvailable // should then exit execution
	}

	// get autosacling groups in the region
	params := &autoscaling.DescribeAutoScalingGroupsInput{}
	autoScalingGroups, err := e.DescribeAutoScalingGroups(params)

	// if error in the returned value, interface custom error
	if err != nil {
		return nil, ErrEC2AutoScalingGroups
	}

	// find the instance auto scaling group
	instanceAutoScalingGroup, err := findInstanceAutoScalingGroup(&instanceIdentity, autoScalingGroups.AutoScalingGroups)

	// check for error
	if err != nil {
		return nil, err
	}

	// get peers instance ids
	instanceIds := make([]*string, 0)
	for _, instance := range instanceAutoScalingGroup.Instances {
		instanceIds = append(instanceIds, instance.InstanceId)
	}

	// peerInstanceIds := filterPeersInstanceIds(instanceIdentity.InstanceID, instanceAutoScalingGroup.Instances)
	peerInstances, err := e.DescribeInstances(&ec2.DescribeInstancesInput{
		InstanceIds: instanceIds,
	})

	// check for error
	if err != nil {
		return nil, err
	}

	// construct return
	instances := make([]*ec2.Instance, 0)
	for _, reservation := range peerInstances.Reservations {
		instances = append(instances, reservation.Instances...)
	}

	// return found peerInstances
	return instances, nil
}

// find asg for instance
func findInstanceAutoScalingGroup(identity *ec2metadata.EC2InstanceIdentityDocument, asgs []*autoscaling.Group) (*autoscaling.Group, error) {
	for _, asg := range asgs {
		// search instances of the asg
		for _, instance := range asg.Instances {
			// check if instance is in the asg
			if *instance.InstanceId == identity.InstanceID {
				// return asg
				return asg, nil
			}
		}
	}

	// return nothing if not found
	return nil, ErrEC2InstanceNotInAutoScalingGroup
}

// filterPeersInstanceIds peer instance ids
func filterPeersInstanceIds(instanceID string, instances []*autoscaling.Instance) []*string {
	ids := make([]*string, 0)
	for _, instance := range instances {
		if instanceID != *instance.InstanceId {
			ids = append(ids, instance.InstanceId)
		}
	}
	return ids
}
