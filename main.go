package main

import (
	"flag"
	"log"

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

	notificationChannel, err := core.InitialiseCloudwatchLogging(sess)
	if err != nil {
		log.Printf("Can't create notification channel %s", err)
	}

	core.RestoreInstances(sess, notificationChannel)

	core.ThrowBanana(sess, configuration.Terminate, notificationChannel)

	core.PokeWithAStick(sess, configuration.Stop, notificationChannel)

}
