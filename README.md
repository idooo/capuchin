# Capuchin

Small Chaos Monkey

- Can terminate/stop instances in autoscaling groups
- Sends notifications to Cloudwatch Logs Streams (capuchin-log-group -> capuchin-log-stream), creates it if needed

## install

```
go get -u github.com/aws/aws-sdk-go
```

