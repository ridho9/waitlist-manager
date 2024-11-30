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

const KeyQueueNumber = "queue-number"
const KeyQueueList = "queue-list"

func incrQueueNumber(ctx context.Context) (int64, error) {
	cmd := vk.B().Incr().Key(KeyQueueNumber).Build()
	resp := vk.Client().Do(ctx, cmd)
	if resp.Error() != nil {
		return 0, resp.Error()
	}
	return resp.AsInt64()
}

func pushQueueNumber(ctx context.Context, queueNumber int64) error {
	cmd := vk.B().Rpush().Key(KeyQueueList).Element(fmt.Sprint(queueNumber)).Build()
	resp := vk.Client().Do(ctx, cmd)
	if resp.Error() != nil {
		return resp.Error()
	}
	return nil
}

func GetQueueList(ctx context.Context) ([]int64, error) {
	cmd := vk.B().Lrange().Key(KeyQueueList).Start(0).Stop(-1).Build()
	resp := vk.Client().Do(ctx, cmd)
	if resp.Error() != nil {
		return nil, resp.Error()
	}
	return resp.AsIntSlice()
}

func GetLastQueueNumber(ctx context.Context) (int64, error) {
	cmd := vk.B().Get().Key(KeyQueueNumber).Build()
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

func QueueCheckIn(ctx context.Context, queueId string) error {
	queueWriteLock.Lock()
	defer queueWriteLock.Unlock()

	queue, err := GetQueueInfo(ctx, queueId)
	if err != nil {
		return err
	}
	if queue == nil {
		return fmt.Errorf("queue not exists")
	}
	if queue.CheckedIn {
		return fmt.Errorf("queue already checked in")
	}
	readyQueueId, _ := GetReadyQueue(ctx)
	if queueId != readyQueueId {
		return fmt.Errorf("this queue is not ready yet")
	}

	queueChairElement := []string{}
	for i := queue.Number; i > 0; i-- {
		queueChairElement = append(queueChairElement, fmt.Sprintf("queue:%s:%d", queueId, i))
	}

	resp := vk.Client().DoMulti(ctx,
		vk.B().Multi().Build(),

		// push clients to chair
		vk.B().Rpush().Key(KeyChairList).Element(queueChairElement...).Build(),

		// set queue checked-in 1
		vk.B().Hset().Key(fmt.Sprintf("party:%s", queueId)).
			FieldValue().
			FieldValue("checked-in", "1").
			Build(),

		vk.B().Lrem().Key(KeyQueueList).Count(1).Element(queueId).Build(),
		vk.B().Set().Key("queue-ready").Value("").Build(),

		vk.B().Exec().Build(),
	)

	for _, r := range resp {
		if r.Error() != nil {
			return fmt.Errorf("transact error: %w", r.Error())
		}
	}

	return nil
}
