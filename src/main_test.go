package main_test

import (
	"bytes"
	"context"
	"log"
	"os"
	"testing"

	notificator "github.com/pippiio/notificator"
	"github.com/pippiio/notificator/handler/slacknotificator"

	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/slack-go/slack"
	"github.com/stretchr/testify/assert"
)

type mockSlackClient struct {
	PostMessageFunc  func(channelID string, options ...slack.MsgOption) (string, string, error)
	UploadFileV2Func func(params slack.UploadFileV2Parameters) (*slack.FileSummary, error)
}

func (m *mockSlackClient) PostMessage(channelID string, options ...slack.MsgOption) (string, string, error) {
	return m.PostMessageFunc(channelID, options...)
}

func (m *mockSlackClient) UploadFileV2(params slack.UploadFileV2Parameters) (*slack.FileSummary, error) {
	return m.UploadFileV2Func(params)
}

var (
	base64Metadata = "ewogICAgImFsbF9maWxlcyI6IFtdLAogICAgImV4aXN0aW5nX2ZpbGVzIjogWwogICAgICAgIHsKICAgICAgICAgICAgInBhdGgiOiAidGVzdC55YW1sIiwKICAgICAgICAgICAgInNpemUiOiAxMCwKICAgICAgICAgICAgImV0YWciOiAiXCJmZGM5ZWU1YmViOGRhYTMyYzE2ZWE0ODUxZTUwY2QwZFwiIiwKICAgICAgICAgICAgImRhdGUiOiAiMjAyNC0wNS0zMFQxMjo0MTo1NS4zMjJaIgogICAgICAgIH0sCiAgICAgICAgewogICAgICAgICAgICAicGF0aCI6ICJ0ZXN0L3Rlc3QxMC55YW1sIiwKICAgICAgICAgICAgInNpemUiOiAxMSwKICAgICAgICAgICAgImV0YWciOiAiXCI2MjZlZDhiZWFlNjI3MDhmZTI1ZWQ5ZjhjMjFmNGUyNlwiIiwKICAgICAgICAgICAgImRhdGUiOiAiMjAyNC0wNi0wNFQwOToyMTozNi4zNTlaIgogICAgICAgIH0sCiAgICAgICAgewogICAgICAgICAgICAicGF0aCI6ICJ0ZXN0MS55YW1sIiwKICAgICAgICAgICAgInNpemUiOiAxMSwKICAgICAgICAgICAgImV0YWciOiAiXCI2MjZlZDhiZWFlNjI3MDhmZTI1ZWQ5ZjhjMjFmNGUyNlwiIiwKICAgICAgICAgICAgImRhdGUiOiAiMjAyNC0wNS0zMFQxNDoxNDozOC43MjRaIgogICAgICAgIH0sCiAgICAgICAgewogICAgICAgICAgICAicGF0aCI6ICJ0ZXN0MTAueWFtbCIsCiAgICAgICAgICAgICJzaXplIjogMTEsCiAgICAgICAgICAgICJldGFnIjogIlwiNjI2ZWQ4YmVhZTYyNzA4ZmUyNWVkOWY4YzIxZjRlMjZcIiIsCiAgICAgICAgICAgICJkYXRlIjogIjIwMjQtMDYtMDRUMDk6MjE6MzAuNjM4WiIKICAgICAgICB9LAogICAgICAgIHsKICAgICAgICAgICAgInBhdGgiOiAidGVzdDIueWFtbCIsCiAgICAgICAgICAgICJzaXplIjogMTEsCiAgICAgICAgICAgICJldGFnIjogIlwiNjI2ZWQ4YmVhZTYyNzA4ZmUyNWVkOWY4YzIxZjRlMjZcIiIsCiAgICAgICAgICAgICJkYXRlIjogIjIwMjQtMDUtMzBUMTI6NTU6NDMuOTUyWiIKICAgICAgICB9LAogICAgICAgIHsKICAgICAgICAgICAgInBhdGgiOiAidGVzdDQueWFtbCIsCiAgICAgICAgICAgICJzaXplIjogMTEsCiAgICAgICAgICAgICJldGFnIjogIlwiNjI2ZWQ4YmVhZTYyNzA4ZmUyNWVkOWY4YzIxZjRlMjZcIiIsCiAgICAgICAgICAgICJkYXRlIjogIjIwMjQtMDUtMzFUMDY6MjU6MDIuNzk1WiIKICAgICAgICB9LAogICAgICAgIHsKICAgICAgICAgICAgInBhdGgiOiAidGVzdDcueWFtbCIsCiAgICAgICAgICAgICJzaXplIjogMTEsCiAgICAgICAgICAgICJldGFnIjogIlwiNjI2ZWQ4YmVhZTYyNzA4ZmUyNWVkOWY4YzIxZjRlMjZcIiIsCiAgICAgICAgICAgICJkYXRlIjogIjIwMjQtMDYtMDRUMDk6MjA6NTYuMDAwWiIKICAgICAgICB9CiAgICBdLAogICAgIm5ld19maWxlcyI6IFtdLAogICAgImNoYW5nZWRfZmlsZXMiOiBbXSwKICAgICJlcnJvcnMiOiBbCiAgICAgICAgewogICAgICAgICAgICAicGF0aCI6ICJ0ZXN0LnlhbWwiLAogICAgICAgICAgICAiZXRhZyI6ICJcImZkYzllZTViZWI4ZGFhMzJjMTZlYTQ4NTFlNTBjZDBkXCIiCiAgICAgICAgfSwKICAgICAgICB7CiAgICAgICAgICAgICJwYXRoIjogInRlc3QvdGVzdDEwLnlhbWwiLAogICAgICAgICAgICAiZXRhZyI6ICJcIjYyNmVkOGJlYWU2MjcwOGZlMjVlZDlmOGMyMWY0ZTI2XCIiCiAgICAgICAgfSwKICAgICAgICB7CiAgICAgICAgICAgICJwYXRoIjogInRlc3QxLnlhbWwiLAogICAgICAgICAgICAiZXRhZyI6ICJcIjYyNmVkOGJlYWU2MjcwOGZlMjVlZDlmOGMyMWY0ZTI2XCIiCiAgICAgICAgfSwKICAgICAgICB7CiAgICAgICAgICAgICJwYXRoIjogInRlc3QxMC55YW1sIiwKICAgICAgICAgICAgImV0YWciOiAiXCI2MjZlZDhiZWFlNjI3MDhmZTI1ZWQ5ZjhjMjFmNGUyNlwiIgogICAgICAgIH0sCiAgICAgICAgewogICAgICAgICAgICAicGF0aCI6ICJ0ZXN0Mi55YW1sIiwKICAgICAgICAgICAgImV0YWciOiAiXCI2MjZlZDhiZWFlNjI3MDhmZTI1ZWQ5ZjhjMjFmNGUyNlwiIgogICAgICAgIH0sCiAgICAgICAgewogICAgICAgICAgICAicGF0aCI6ICJ0ZXN0NC55YW1sIiwKICAgICAgICAgICAgImV0YWciOiAiXCI2MjZlZDhiZWFlNjI3MDhmZTI1ZWQ5ZjhjMjFmNGUyNlwiIgogICAgICAgIH0sCiAgICAgICAgewogICAgICAgICAgICAicGF0aCI6ICJ0ZXN0Ny55YW1sIiwKICAgICAgICAgICAgImV0YWciOiAiXCI2MjZlZDhiZWFlNjI3MDhmZTI1ZWQ5ZjhjMjFmNGUyNlwiIgogICAgICAgIH0KICAgIF0sCiAgICAidG90YWxfYWxsX2ZpbGVzIjogMCwKICAgICJ0b3RhbF9leGlzdGluZ19maWxlcyI6IDcsCiAgICAidG90YWxfbmV3X2ZpbGVzIjogMCwKICAgICJ0b3RhbF9jaGFuZ2VkX2ZpbGVzIjogMCwKICAgICJ0b3RhbF9lcnJvcnMiOiA3Cn0K"
)

func TestHandler_Slack_Success(t *testing.T) {
	var buf bytes.Buffer
	log.SetOutput(&buf)
	defer func() {
		log.SetOutput(os.Stderr)
	}()

	slackClient := &mockSlackClient{}

	slackClient.PostMessageFunc = func(channelID string, options ...slack.MsgOption) (string, string, error) {
		return "1234", "4321", nil
	}

	slackClient.UploadFileV2Func = func(params slack.UploadFileV2Parameters) (*slack.FileSummary, error) {
		return &slack.FileSummary{}, nil
	}

	ctx := slacknotificator.WithClient(context.Background(), &slacknotificator.Client{SlackClient: slackClient})

	resp, err := notificator.Handler(ctx, events.SQSEvent{
		Records: []events.SQSMessage{
			{
				Body: ":white_check_mark: *Info in bucket `test`* \n *Message:* Backup of bucket test processed successfully",
				MessageAttributes: map[string]events.SQSMessageAttribute{
					"channel": {
						StringValue: aws.String("slack"),
					},
					"file-metadata": {
						StringValue: &base64Metadata,
					},
				},
			},
			{
				Body: ":warning: *Error in bucket `test-error`* \n *Message:* Metadata file contains errors",
				MessageAttributes: map[string]events.SQSMessageAttribute{
					"channel": {
						StringValue: aws.String("slack"),
					},
					"file-metadata": {
						StringValue: &base64Metadata,
					},
				},
			},
			{
				Body: ":no_entry: *Fatal in bucket `test-fatal`* \n *Message:* No metadata file present for today.",
				MessageAttributes: map[string]events.SQSMessageAttribute{
					"channel": {
						StringValue: aws.String("slack"),
					},
				},
			},
		},
	})

	t.Log(buf.String())

	assert.NoError(t, err)
	assert.NotNil(t, resp)
	assert.Regexp(t, "(?s).*Slack message sent(?s).*", buf.String())
}

func TestHandler_No_MessageAttribute(t *testing.T) {
	var buf bytes.Buffer
	log.SetOutput(&buf)
	defer func() {
		log.SetOutput(os.Stderr)
	}()

	resp, err := notificator.Handler(context.Background(), events.SQSEvent{
		Records: []events.SQSMessage{
			{
				Body: ":white_check_mark: *Info in bucket `test`* \n *Message:* Backup of bucket test processed successfully",
			},
		},
	})

	t.Log(buf.String())

	assert.NoError(t, err)
	assert.NotNil(t, resp)
	assert.Regexp(t, "(?s).*No message attributes(?s).*", buf.String())
}

func TestHandler_No_Channel(t *testing.T) {
	var buf bytes.Buffer
	log.SetOutput(&buf)
	defer func() {
		log.SetOutput(os.Stderr)
	}()

	resp, err := notificator.Handler(context.Background(), events.SQSEvent{
		Records: []events.SQSMessage{
			{
				Body: ":white_check_mark: *Info in bucket `test`* \n *Message:* Backup of bucket test processed successfully",
				MessageAttributes: map[string]events.SQSMessageAttribute{
					"file-metadata": {
						StringValue: aws.String("test"),
					},
				},
			},
		},
	})

	t.Log(buf.String())

	assert.NoError(t, err)
	assert.NotNil(t, resp)
	assert.Regexp(t, "(?s).*No channel attribute(?s).*", buf.String())
}

func TestHandler_Unknown_Channel(t *testing.T) {
	var buf bytes.Buffer
	log.SetOutput(&buf)
	defer func() {
		log.SetOutput(os.Stderr)
	}()

	resp, err := notificator.Handler(context.Background(), events.SQSEvent{
		Records: []events.SQSMessage{
			{
				Body: ":white_check_mark: *Info in bucket `test`* \n *Message:* Backup of bucket test processed successfully",
				MessageAttributes: map[string]events.SQSMessageAttribute{
					"channel": {
						StringValue: aws.String("unknown"),
					},
					"file-metadata": {
						StringValue: aws.String("test"),
					},
				},
			},
		},
	})

	t.Log(buf.String())

	assert.NoError(t, err)
	assert.NotNil(t, resp)
	assert.Regexp(t, "(?s).*Unknown channel(?s).*", buf.String())
}
