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
		*notificationChannel <- err.Error()
		return err
	}

	*notificationChannel <- fmt.Sprintf("Throwing banana at %s", *pickedInstance.InstanceId)

	err = TerminateInstance(currentSession, pickedInstance, config.Tag)
	if err != nil {
		*notificationChannel <- err.Error()
		return err
	}

	*notificationChannel <- fmt.Sprintf("Instance %s has been killed by a banana", *pickedInstance.InstanceId)

	return nil
}

// PokeWithAStick - stops an eligible instance
func PokeWithAStick(currentSession *session.Session, config *StateControlConfiguration) error {

	notificationChannel := GetNotificationChannel(currentSession)
	*notificationChannel <- fmt.Sprintf("Attempting to poke something %v with a stick", config.Instances)

	pickedInstance, err := PickInstance(currentSession, nil, config.Instances)
	if err != nil {
		*notificationChannel <- err.Error()
		return err
	}

	*notificationChannel <- fmt.Sprintf("Poking %s with a stick", *pickedInstance.InstanceId)

	err = StopInstance(currentSession, pickedInstance, config.Tag)
	if err != nil {
		*notificationChannel <- err.Error()
		return err
	}

	*notificationChannel <- fmt.Sprintf("Instance %s has been stopped by a stick", *pickedInstance.InstanceId)

	return nil
}

// RestoreInstances - starts previously stopped instances
func RestoreInstances(currentSession *session.Session, config *StateControlConfiguration) error {

	notificationChannel := GetNotificationChannel(currentSession)
	*notificationChannel <- fmt.Sprintf("Attempting to restore %v", config.Instances)

	pickedInstances, err := PickStoppedInstances(currentSession, config.Instances)
	if err != nil {
		*notificationChannel <- err.Error()
		return err
	}

	var instanceIds []*string
	for _, instance := range pickedInstances {
		instanceIds = append(instanceIds, instance.InstanceId)
	}
	*notificationChannel <- fmt.Sprintf("Restoring instances %v ...", instanceIds)

	err = StartInstances(currentSession, pickedInstances, config.Tag)
	if err != nil {
		*notificationChannel <- err.Error()
		return err
	}

	return nil
}
