package core

import (
	"errors"
	"log"
	"time"

	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/cloudwatchlogs"
)

func getLogGroup(service *cloudwatchlogs.CloudWatchLogs, logGroupName *string) (*cloudwatchlogs.LogGroup, error) {
	search := &cloudwatchlogs.DescribeLogGroupsInput{
		LogGroupNamePrefix: logGroupName,
	}
	describeLogGroupsOutput, err := service.DescribeLogGroups(search)
	if err != nil {
		return nil, errors.New("Can't get cloudwatch log groups")
	}

	if len(describeLogGroupsOutput.LogGroups) != 0 {
		return describeLogGroupsOutput.LogGroups[0], nil
	}

	return nil, nil
}

func createLogGroup(service *cloudwatchlogs.CloudWatchLogs) (*string, error) {

	logGroupName := "capuchin-log-group"
	logGroupRole := "capuchin"
	logGroupTags := map[string]*string{
		"Role": &logGroupRole,
	}

	logGroup, err := getLogGroup(service, &logGroupName)
	if err != nil {
		return nil, err
	}
	if logGroup != nil {
		return logGroup.LogGroupName, nil
	}

	log.Printf("Creating new Cloudwatch log group: %s", logGroupName)
	input := &cloudwatchlogs.CreateLogGroupInput{
		LogGroupName: &logGroupName,
		Tags:         logGroupTags,
	}
	_, err = service.CreateLogGroup(input)
	if err != nil {
		return nil, errors.New("Can't create cloudwatch log group")
	}

	return &logGroupName, nil
}

func getLogStream(service *cloudwatchlogs.CloudWatchLogs, logGroupName *string, logStreamName *string) (*cloudwatchlogs.LogStream, error) {
	search := &cloudwatchlogs.DescribeLogStreamsInput{
		LogGroupName:        logGroupName,
		LogStreamNamePrefix: logStreamName,
	}
	describeLogStreamsOutput, err := service.DescribeLogStreams(search)
	if err != nil {
		return nil, errors.New("Can't get cloudwatch log streams")
	}

	if len(describeLogStreamsOutput.LogStreams) != 0 {
		return describeLogStreamsOutput.LogStreams[0], nil
	}

	return nil, nil
}

func createLogStream(service *cloudwatchlogs.CloudWatchLogs, logGroupName *string) (*string, error) {

	logStreamName := "capuchin-log-stream"
	logStream, err := getLogStream(service, logGroupName, &logStreamName)
	if err != nil {
		return nil, err
	}
	if logStream != nil {
		return logStream.LogStreamName, nil
	}

	log.Printf("Creating new Cloudwatch log stream: %s", logStreamName)
	input := &cloudwatchlogs.CreateLogStreamInput{
		LogGroupName:  logGroupName,
		LogStreamName: &logStreamName,
	}
	_, err = service.CreateLogStream(input)
	if err != nil {
		return nil, errors.New("Can't create cloudwatch log stream")
	}
	return &logStreamName, nil
}

func sendCloudwatchEvent(service *cloudwatchlogs.CloudWatchLogs, logGroupName *string, logStreamName *string, message *string) {
	logStream, err := getLogStream(service, logGroupName, logStreamName)
	if err != nil {
		log.Printf("Error retrieving Cloudwatch log stream: %s", err)
		return
	}
	currentTimestamp := int64(time.Now().UnixNano() / int64(time.Millisecond))
	logEventInput := &cloudwatchlogs.PutLogEventsInput{
		LogGroupName:  logGroupName,
		LogStreamName: logStreamName,
		SequenceToken: logStream.UploadSequenceToken,
		LogEvents: []*cloudwatchlogs.InputLogEvent{
			{
				Message:   message,
				Timestamp: &currentTimestamp,
			},
		},
	}

	_, err = service.PutLogEvents(logEventInput)
	if err != nil {
		log.Printf("Error publishing event to Cloudwatch: %s", err)
	}
}

func listenForChannel(service *cloudwatchlogs.CloudWatchLogs, logGroupName *string, logStreamName *string, notificationChannel chan string) {

	for message := range notificationChannel {
		log.Println(message)
		sendCloudwatchEvent(service, logGroupName, logStreamName, &message)
	}
}

func InitialiseCloudwatchLogging(currentSession *session.Session) (*chan string, error) {
	service := cloudwatchlogs.New(currentSession)
	logGroupName, err := createLogGroup(service)
	if err != nil {
		return nil, err
	}
	logStreamName, err := createLogStream(service, logGroupName)
	if err != nil {
		return nil, err
	}

	log.Printf("Using Cloudwatch log group and stream: %s -> %s", *logGroupName, *logStreamName)

	notificationChannel := make(chan string)
	go listenForChannel(service, logGroupName, logStreamName, notificationChannel)

	return &notificationChannel, nil
}
