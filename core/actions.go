package core

import (
	"fmt"

	"github.com/aws/aws-sdk-go/aws/session"
)

// ThrowBanana - throws a banana to terminate eligible instance in Autoscaling group
func ThrowBanana(currentSession *session.Session, config *StateControlConfiguration) error {

	notificationChannel := GetNotificationChannel(currentSession)
	*notificationChannel <- fmt.Sprintf("Attempting to throw banana at something %v -> %v", config.Autoscaling, config.Instances)

	pickedGroup, err := PickAutoscalingGroup(currentSession, config.Autoscaling)
	if err != nil {
		return err
	}
	*notificationChannel <- fmt.Sprintf("Picked group: %s", *pickedGroup.AutoScalingGroupName)

	pickedInstance, err := PickInstance(currentSession, pickedGroup, config.Instances)
	if err != nil {
		return err
	}

	*notificationChannel <- fmt.Sprintf("Throwing banana at %s", *pickedInstance.InstanceId)

	// TODO terminate instance
	// add tag that terminated by a monkey

	return nil
}

// PokeWithAStick - stops an eligible instance
func PokeWithAStick(currentSession *session.Session, config *StateControlConfiguration) error {

	notificationChannel := GetNotificationChannel(currentSession)
	*notificationChannel <- fmt.Sprintf("Attempting to poke something %v with a stick", config.Instances)

	pickedInstance, err := PickInstance(currentSession, nil, config.Instances)
	if err != nil {
		return err
	}

	*notificationChannel <- fmt.Sprintf("Poking %s with a stick", *pickedInstance.InstanceId)

	// TODO stop instance
	// add tag that stopped by a monkey

	return nil
}

// RestoreInstances - starts previously stopped instances
func RestoreInstances(sess *session.Session) {
	// TODO: implement
}
