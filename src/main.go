package main

import (
	"context"
	"log"

	"github.com/pippiio/notificator/handler/slacknotificator"

	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
)

func Handler(ctx context.Context, event events.SQSEvent) (events.APIGatewayProxyResponse, error) {
	log.SetFlags(log.Ldate | log.Ltime | log.Lshortfile | log.Lmsgprefix)

	flags, err := InitFlags()
	if err != nil {
		log.Fatal(err)
	}

	for _, record := range event.Records {
		if record.MessageAttributes == nil {
			log.SetPrefix("ERROR: ")
			log.Printf("No message attributes, skipping, record: %s", record.Body)
			continue
		}

		value, ok := record.MessageAttributes["channel"]
		if !ok {
			log.SetPrefix("ERROR: ")
			log.Printf("No channel attribute, skipping, record: %s", record.Body)
			continue
		}

		switch *value.StringValue {
		case "slack":
			client, err := slacknotificator.New(ctx, flags.Token, flags.ChannelID)
			if err != nil {
				log.Fatal(err)
			}
			ctx = slacknotificator.WithClient(ctx, client) // singleton

			err = SlackHandler(client, flags, record)
			if err != nil {
				log.SetPrefix("ERROR: ")
				log.Print(err.Error())
				continue
			}

			log.SetPrefix("INFO: ")
			log.Printf("Slack message sent, record: %s", record.Body)
		default:
			log.SetPrefix("ERROR: ")
			log.Printf("Unknown channel, skipping, record: %s", record.Body)
			continue
		}
	}

	response := events.APIGatewayProxyResponse{
		StatusCode: 200,
	}
	return response, nil
}

func main() {
	lambda.Start(Handler)
}
