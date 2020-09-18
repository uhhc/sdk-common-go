package httpclient

import (
	"time"

	"github.com/go-resty/resty/v2"
	"github.com/spf13/viper"
)

func init() {
	viper.AutomaticEnv()
}

// NewClient to return a Resty http client
func NewClient() *resty.Client {
	// Create a Resty Client
	client := resty.New()

	// Set client timeout
	timeout := viper.GetInt64("HTTP_CLIENT_TIMEOUT")
	if timeout == 0 {
		timeout = 10
	}
	client.SetTimeout(time.Duration(timeout) * time.Second)

	return client
}
