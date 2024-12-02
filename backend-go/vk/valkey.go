package vk

import (
	"backend-go/env"

	"github.com/valkey-io/valkey-go"
)

var client valkey.Client

func init() {
	var err error
	client, err = valkey.NewClient(valkey.ClientOption{InitAddress: []string{env.VALKEY_ADDRESS}})
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
