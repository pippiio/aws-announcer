package main

import (
	"github.com/Netflix/go-env"
	"github.com/joho/godotenv"
)

type Flags struct {
	Token     string `env:"SLACK_TOKEN"`
	ChannelID string `env:"SLACK_CHANNELID"`
}

// InitFlags initializes the Flags struct by loading environment variables and unmarshaling them into the struct.
//
// It returns a pointer to the initialized Flags struct and an error if any occurred during the initialization process.
func InitFlags() (*Flags, error) {
	_ = godotenv.Load()

	var flags Flags
	_, err := env.UnmarshalFromEnviron(&flags)
	if err != nil {
		return nil, err
	}
	return &flags, nil
}
