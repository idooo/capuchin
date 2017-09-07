package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"log"
	"math/rand"
	"os"
	"time"

	"github.com/aws/aws-sdk-go/aws/awserr"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/autoscaling"
	"github.com/aws/aws-sdk-go/service/ec2"
)

type StateControlConfiguration struct {
	Instances   *map[string]string
	Autoscaling *map[string]string
}

type Configuration struct {
	Terminate *StateControlConfiguration
	Stop      *StateControlConfiguration
}

// Reads configuration file from the specified location and
// applies the default values if needed
func readConfig(filename string) Configuration {
	file, _ := os.Open(filename)
	decoder := json.NewDecoder(file)
	configuration := Configuration{}
	err := decoder.Decode(&configuration)
	if err != nil {
		log.Printf("Can't read configuration file %s : %s", filename, err)
	}
	return configuration
}

// whatever
func RetriveEligibleAutoscalingGroup(currentSession *session.Session, asgTags *map[string]string) []*autoscaling.Group {
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
			// Print the error, cast err to awserr.Error to get the Code and
			// Message from an error.
			log.Fatal(err.Error())
		}
	}

	var eligibleAutoScalingGroups []*autoscaling.Group
	for _, group := range result.AutoScalingGroups {
		// fmt.Println(element.Tags)

		isEligible := true
		for requiredKey, requiredValue := range *asgTags {
			isTagFound := false
			for _, tag := range group.Tags {
				if requiredKey == *tag.Key && requiredValue == *tag.Value {
					isTagFound = true
					break
				}
			}
			isEligible = isEligible && isTagFound
		}
		if isEligible {
			eligibleAutoScalingGroups = append(eligibleAutoScalingGroups, group)
		}
		// element is the element from someSlice for where we are
	}
	log.Printf("Found %d of %d eligible autoscaling groups", len(eligibleAutoScalingGroups), len(result.AutoScalingGroups))
	return eligibleAutoScalingGroups
}

func RetrieveEligibleInstances(currentSession *session.Session, group *autoscaling.Group) {
	svc := ec2.New(currentSession)

	var instanceIds []*string
	for _, instance := range group.Instances {
		instanceIds = append(instanceIds, instance.InstanceId)
	}

	input := &ec2.DescribeInstancesInput{
		InstanceIds: instanceIds,
	}

	result, err := svc.DescribeInstances(input)
	if err != nil {
		fmt.Println("Error", err)
	} else {
		fmt.Println("Success", result)
	}
}

func main() {

	configPathPtr := flag.String(
		"config",
		"./config/config.json",
		"path to a configuration file")

	flag.Parse()
	configuration := readConfig(*configPathPtr)

	sess := session.Must(session.NewSessionWithOptions(session.Options{
		SharedConfigState: session.SharedConfigEnable,
	}))

	eligibleAutoScalingGroups := RetriveEligibleAutoscalingGroup(sess, configuration.Stop.Autoscaling)

	rand.Seed(time.Now().Unix())
	pickedGroup := eligibleAutoScalingGroups[rand.Intn(len(eligibleAutoScalingGroups))]

	log.Printf("Picked group: %s", *pickedGroup.AutoScalingGroupName)

	RetrieveEligibleInstances(sess, pickedGroup)

	// // Create new EC2 client
	// ec2Svc := ec2.New(sess)

	// result, err := ec2Svc.DescribeInstances(nil)
	// if err != nil {
	// 	fmt.Println("Error", err)
	// } else {
	// 	fmt.Println("Success", result)
	// }
}
