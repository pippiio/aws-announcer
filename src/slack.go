package main

import (
	"encoding/base64"
	"strings"

	"github.com/pippiio/notificator/handler/slacknotificator"

	"github.com/aws/aws-lambda-go/events"
)

func SlackHandler(client *slacknotificator.Client, flags *Flags, record events.SQSMessage) error {
	timestamp, err := client.SendSlackMessage(record.Body, "")
	if err != nil {
		return err
	}

	for key, attribute := range record.MessageAttributes {
		if strings.HasPrefix(key, "file-") {
			fileName, _ := strings.CutPrefix(key, "file-")

			content, err := base64.StdEncoding.DecodeString(*attribute.StringValue)
			if err != nil {
				return err
			}

			err = client.SendSlackFile(string(content), fileName, timestamp)
			if err != nil {
				return err
			}
		}
	}

	return nil
}
