package main

import (
	"flag"

	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/idooo/capuchin/core"
)

func main() {

	configPathPtr := flag.String(
		"config",
		"./config/config.json",
		"path to a configuration file")

	flag.Parse()
	configuration := core.ReadConfigiguration(*configPathPtr)

	sess := session.Must(session.NewSessionWithOptions(session.Options{
		SharedConfigState: session.SharedConfigEnable,
	}))

	// Create notification channel
	notificationChannel := core.GetNotificationChannel(sess)
	*notificationChannel <- "Capuchin has been awakened..."

	// Do what you have to do
	if configuration.Restore != nil {
		core.RestoreInstances(sess, configuration.Restore)
	}

	if configuration.Terminate != nil {
		core.ThrowBanana(sess, configuration.Terminate)
	}

	if configuration.Stop != nil {
		core.PokeWithAStick(sess, configuration.Stop)
	}
}
