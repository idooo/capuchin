package core

import (
	"log"

	"github.com/aws/aws-sdk-go/aws/session"
)

// ThrowBanana - throws a banana to terminate eligible instance in Autoscaling group
func ThrowBanana(sess *session.Session, config *StateControlConfiguration) error {
	log.Printf("Attempting to throw banana at something %v -> %v", config.Autoscaling, config.Instances)

	pickedGroup, err := PickAutoscalingGroup(sess, config.Autoscaling)
	if err != nil {
		return err
	}
	log.Printf("Picked group: %s", *pickedGroup.AutoScalingGroupName)

	pickedInstance, err := PickInstance(sess, pickedGroup, config.Instances)
	if err != nil {
		return err
	}

	log.Printf("Throwing banana at %s", *pickedInstance.InstanceId)

	// TODO terminate instance
	// add tag that terminated by a monkey

	return nil
}

// PokeWithAStick - stops an eligible instance
func PokeWithAStick(sess *session.Session, config *StateControlConfiguration) error {

	log.Printf("Attempting to poke something %v with a stick", config.Instances)
	pickedInstance, err := PickInstance(sess, nil, config.Instances)
	if err != nil {
		return err
	}

	log.Printf("Poking %s with a stick", *pickedInstance.InstanceId)

	// TODO stop instance
	// add tag that stopped by a monkey

	return nil
}

// RestoreInstances - starts previously stopped instances
func RestoreInstances(sess *session.Session) {
	// TODO: implement
}
