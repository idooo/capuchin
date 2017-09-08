package core

func sendNotificationToHipChat() {

}

func sendLogsToSumoLogic() {

}

// Notify - sends notifications and logs
func Notify(message string) {
	sendNotificationToHipChat()
	sendLogsToSumoLogic()
}
