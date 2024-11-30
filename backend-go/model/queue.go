package model

import (
	"backend-go/valkey"
	"context"
	"fmt"
	"strconv"
	"sync"
)

type QueueInfo struct {
	Name   string `json:"name"`
	Number int64  `json:"number"`
}

var queueWriteLock sync.Mutex

func incrQueueNumber(ctx context.Context) (int64, error) {
	cmd := valkey.B().Incr().Key("queue-number").Build()
	resp := valkey.Client().Do(ctx, cmd)
	if resp.Error() != nil {
		return 0, resp.Error()
	}
	return resp.AsInt64()
}

func pushQueueNumber(ctx context.Context, queueNumber int64) error {
	cmd := valkey.B().Rpush().Key("queue-list").Element(fmt.Sprint(queueNumber)).Build()
	resp := valkey.Client().Do(ctx, cmd)
	if resp.Error() != nil {
		return resp.Error()
	}
	return nil
}

func GetQueueList(ctx context.Context) ([]int64, error) {
	cmd := valkey.B().Lrange().Key("queue-list").Start(0).Stop(-1).Build()
	resp := valkey.Client().Do(ctx, cmd)
	if resp.Error() != nil {
		return nil, resp.Error()
	}
	return resp.AsIntSlice()
}

func GetLastQueueNumber(ctx context.Context) (int64, error) {
	cmd := valkey.B().Get().Key("queue-number").Build()
	resp := valkey.Client().Do(ctx, cmd)
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

	cmd := valkey.B().Hset().
		Key(fmt.Sprintf("party:%d", newQueueNumber)).
		FieldValue().
		FieldValue("name", partyName).
		FieldValue("number", fmt.Sprint(partyNumber)).
		Build()

	resp := valkey.Client().Do(ctx, cmd)
	if resp.Error() != nil {
		return 0, resp.Error()
	}

	return newQueueNumber, nil
}

func GetQueueInfo(ctx context.Context, queueId string) (*QueueInfo, error) {
	cmd := valkey.B().Hgetall().Key(fmt.Sprintf("party:%s", queueId)).Build()
	resp := valkey.Client().Do(ctx, cmd)
	if resp.Error() != nil {
		return nil, resp.Error()
	}

	arr, _ := resp.ToArray()
	fmt.Println(arr)

	vals, err := resp.AsStrMap()
	if err != nil {
		return nil, err
	}
	if len(vals) == 0 {
		return nil, nil
	}

	name := vals["name"]
	number, _ := strconv.Atoi(vals["number"])

	return &QueueInfo{Name: name, Number: int64(number)}, nil
}
