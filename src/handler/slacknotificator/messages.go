package slacknotificator

import (
	"bytes"

	"github.com/slack-go/slack"
)

func (s *Client) SendSlackMessage(message string, timestamp string) (string, error) {
	_, timestamp, err := s.SlackClient.PostMessage(s.ChannelID, slack.MsgOptionText(message, false), slack.MsgOptionTS(timestamp))
	if err != nil {
		return "", err
	}

	return timestamp, nil
}

func (s *Client) SendSlackFile(fileContent, fileName, threadTimestamp string) error {
	fileUploadParams := slack.UploadFileV2Parameters{
		Filename:        fileName,
		FileSize:        len(fileContent),
		Reader:          bytes.NewReader([]byte(fileContent)),
		Channel:         s.ChannelID,
		ThreadTimestamp: threadTimestamp,
	}

	_, err := s.SlackClient.UploadFileV2(fileUploadParams)
	if err != nil {
		return err
	}

	return nil
}
