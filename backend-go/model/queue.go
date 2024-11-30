package model

import (
	"backend-go/vk"
	"context"
	"fmt"
	"strconv"
	"sync"

	"github.com/valkey-io/valkey-go"
)

type QueueInfo struct {
	Name      string `json:"name"`
	Number    int64  `json:"number"`
	CheckedIn bool   `json:"checked_in"`
}

var queueWriteLock sync.Mutex

func incrQueueNumber(ctx context.Context) (int64, error) {
	cmd := vk.B().Incr().Key("queue-number").Build()
	resp := vk.Client().Do(ctx, cmd)
	if resp.Error() != nil {
		return 0, resp.Error()
	}
	return resp.AsInt64()
}

func pushQueueNumber(ctx context.Context, queueNumber int64) error {
	cmd := vk.B().Rpush().Key("queue-list").Element(fmt.Sprint(queueNumber)).Build()
	resp := vk.Client().Do(ctx, cmd)
	if resp.Error() != nil {
		return resp.Error()
	}
	return nil
}

func GetQueueList(ctx context.Context) ([]int64, error) {
	cmd := vk.B().Lrange().Key("queue-list").Start(0).Stop(-1).Build()
	resp := vk.Client().Do(ctx, cmd)
	if resp.Error() != nil {
		return nil, resp.Error()
	}
	return resp.AsIntSlice()
}

func GetLastQueueNumber(ctx context.Context) (int64, error) {
	cmd := vk.B().Get().Key("queue-number").Build()
	resp := vk.Client().Do(ctx, cmd)
	if resp.Error() != nil {
		return 0, resp.Error()
	}
	return resp.AsInt64()
}

func AddNewQueue(ctx context.Context, partyName string, partyNumber int) (int64, error) {
	queueWriteLock.Lock()
	defer queueWriteLock.Unlock()

	newQueueNumber, err := incrQueueNumber(ctx)
	if err != nil {
		return 0, err
	}
	err = pushQueueNumber(ctx, newQueueNumber)
	if err != nil {
		return 0, err
	}

	cmd := vk.B().Hset().
		Key(fmt.Sprintf("party:%d", newQueueNumber)).
		FieldValue().
		FieldValue("name", partyName).
		FieldValue("number", fmt.Sprint(partyNumber)).
		FieldValue("checked-in", "0").
		Build()

	resp := vk.Client().Do(ctx, cmd)
	if resp.Error() != nil {
		return 0, resp.Error()
	}

	return newQueueNumber, nil
}

func GetQueueInfo(ctx context.Context, queueId string) (*QueueInfo, error) {
	cmd := vk.B().Hgetall().Key(fmt.Sprintf("party:%s", queueId)).Build()
	resp := vk.Client().Do(ctx, cmd)
	if resp.Error() != nil {
		return nil, resp.Error()
	}

	vals, err := resp.AsStrMap()
	if err != nil {
		return nil, err
	}
	if len(vals) == 0 {
		return nil, nil
	}

	name := vals["name"]
	number, _ := strconv.Atoi(vals["number"])
	checkedIn := vals["checked-in"] == "1"

	return &QueueInfo{Name: name, Number: int64(number), CheckedIn: checkedIn}, nil
}

func GetReadyQueue(ctx context.Context) (string, error) {
	cmd := vk.B().Get().Key("queue-ready").Build()
	resp := vk.Client().Do(ctx, cmd)

	if resp.Error() != nil {
		if resp.Error() == valkey.Nil {
			return "", nil
		}

		return "", resp.Error()
	}

	return resp.ToString()
}

func SetReadyQueue(ctx context.Context, queueId string) error {
	cmd := vk.B().Set().Key("queue-ready").Value(queueId).Build()
	resp := vk.Client().Do(ctx, cmd)
	return resp.Error()
}
