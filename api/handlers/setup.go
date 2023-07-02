package handlers

import (
	"log"

	"github.com/stripe/stripe-go/v74/client"
)

type StripeConfig struct {
	Api *client.API
	Log *log.Logger
}

func NewStripeConfig(key string, l *log.Logger) StripeConfig {
	api := client.New(key, nil)
	return StripeConfig{api, l}
}
