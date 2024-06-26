package slacknotificator

import (
	"context"

	"github.com/slack-go/slack"
)

type SlackClient interface {
	PostMessage(channelID string, options ...slack.MsgOption) (string, string, error)
	UploadFileV2(params slack.UploadFileV2Parameters) (*slack.FileSummary, error)
}

type Client struct {
	SlackClient SlackClient
	ChannelID   string
}

type contextKey string

var (
	clientContextKey contextKey = "slackclient.client"
)

// New creates a new Client instance using the provided context, profile, region, and endpoint.
//
// Parameters:
// - ctx: The context.Context object for cancellation and timeout control.
// - profile: The AWS profile to use.
// - region: The AWS region to use.
// - endpoint: The custom endpoint to use for S3.
//
// Returns:
// - *Client: A pointer to the newly created Client instance.
// - error: An error if the operation fails.
func New(ctx context.Context, token, channelID string) (*Client, error) {
	client, ok := ctx.Value(clientContextKey).(*Client)
	if ok && client != nil {
		return client, nil
	}

	slackClient := slack.New(token)

	return &Client{
		ChannelID:   channelID,
		SlackClient: slackClient,
	}, nil
}

// WithClient updates the context with the provided client information.
//
// - ctx: The context.Context object for the operation.
// - client: A pointer to the Client struct containing client information.
// Returns:
// - context.Context: The updated context with the client information.
func WithClient(ctx context.Context, client *Client) context.Context {
	return context.WithValue(ctx, clientContextKey, client)
}
