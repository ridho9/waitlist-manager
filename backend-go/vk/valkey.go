package vk

import (
	"github.com/valkey-io/valkey-go"
)

var client valkey.Client

func init() {
	var err error
	client, err = valkey.NewClient(valkey.ClientOption{InitAddress: []string{"127.0.0.1:8502"}})
	if err != nil {
		panic(err)
	}
}

func Client() valkey.Client {
	return client
}

func B() valkey.Builder {
	return client.B()
}
