package core

import (
	"errors"
	"log"
	"math/rand"
	"time"

	"github.com/aws/aws-sdk-go/aws/awserr"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/autoscaling"
	"github.com/aws/aws-sdk-go/service/ec2"
)

func isEligibleAutoscalingGroup(group *autoscaling.Group, requiredTags *map[string]string) bool {

	isEligible := true
	for requiredKey, requiredValue := range *requiredTags {
		isTagFound := false
		for _, tag := range group.Tags {
			if requiredKey == *tag.Key && requiredValue == *tag.Value {
				isTagFound = true
				break
			}
		}
		isEligible = isEligible && isTagFound
	}
	return isEligible
}

func isEligibleInstance(instance *ec2.Instance, requiredTags *map[string]string) bool {

	// Only accept instances that have running state
	runningCode := int64(16)
	if instance.State.Code != &runningCode {
		return false
	}

	isEligible := true
	for requiredKey, requiredValue := range *requiredTags {
		isTagFound := false
		for _, tag := range instance.Tags {
			if requiredKey == *tag.Key && requiredValue == *tag.Value {
				isTagFound = true
				break
			}
		}
		isEligible = isEligible && isTagFound
	}
	return isEligible
}

func retriveEligibleAutoscalingGroup(currentSession *session.Session, requiredTags *map[string]string) []*autoscaling.Group {
	svc := autoscaling.New(currentSession)
	input := &autoscaling.DescribeAutoScalingGroupsInput{
		MaxRecords: &[]int64{100}[0],
	}

	result, err := svc.DescribeAutoScalingGroups(input)
	if err != nil {
		if aerr, ok := err.(awserr.Error); ok {
			switch aerr.Code() {
			case autoscaling.ErrCodeInvalidNextToken:
				log.Fatal(autoscaling.ErrCodeInvalidNextToken, aerr.Error())
			case autoscaling.ErrCodeResourceContentionFault:
				log.Fatal(autoscaling.ErrCodeResourceContentionFault, aerr.Error())
			default:
				log.Fatal(aerr.Error())
			}
		} else {
			log.Fatal(err.Error())
		}
	}

	var eligibleAutoScalingGroups []*autoscaling.Group
	for _, group := range result.AutoScalingGroups {
		if isEligibleAutoscalingGroup(group, requiredTags) {
			eligibleAutoScalingGroups = append(eligibleAutoScalingGroups, group)
		}
	}
	log.Printf("Found %d of %d eligible autoscaling groups", len(eligibleAutoScalingGroups), len(result.AutoScalingGroups))
	return eligibleAutoScalingGroups
}

func retrieveEligibleInstances(currentSession *session.Session, group *autoscaling.Group, requiredTags *map[string]string) []*ec2.Instance {
	svc := ec2.New(currentSession)

	input := &ec2.DescribeInstancesInput{}

	if group != nil {
		var instanceIds []*string
		for _, instance := range group.Instances {
			instanceIds = append(instanceIds, instance.InstanceId)
		}

		input = &ec2.DescribeInstancesInput{
			InstanceIds: instanceIds,
		}
	}

	result, err := svc.DescribeInstances(input)
	if err != nil {
		log.Fatal("Error", err)
	}

	var eligibleInstances []*ec2.Instance
	for _, instance := range result.Reservations[0].Instances {

		if isEligibleInstance(instance, requiredTags) {
			eligibleInstances = append(eligibleInstances, instance)
		}
	}
	log.Printf("Found %d of %d eligible instances", len(eligibleInstances), len(result.Reservations[0].Instances))
	return eligibleInstances
}

// PickAutoscalingGroup - picks eligible autoscaling group
func PickAutoscalingGroup(currentSession *session.Session, requiredTags *map[string]string) (*autoscaling.Group, error) {
	rand.Seed(time.Now().Unix())
	eligibleAutoScalingGroups := retriveEligibleAutoscalingGroup(currentSession, requiredTags)

	if len(eligibleAutoScalingGroups) == 0 {
		return nil, errors.New("No eligible Autoscaling Groups found")
	}
	pickedGroup := eligibleAutoScalingGroups[rand.Intn(len(eligibleAutoScalingGroups))]

	return pickedGroup, nil
}

// PickInstance - picks eligible instance
func PickInstance(currentSession *session.Session, autoscalingGroup *autoscaling.Group, requiredTags *map[string]string) (*ec2.Instance, error) {
	rand.Seed(time.Now().Unix())
	eligibleInstnaces := retrieveEligibleInstances(currentSession, autoscalingGroup, requiredTags)

	if len(eligibleInstnaces) == 0 {
		return nil, errors.New("No eligible instances found")
	}
	pickedInstance := eligibleInstnaces[rand.Intn(len(eligibleInstnaces))]

	return pickedInstance, nil
}
