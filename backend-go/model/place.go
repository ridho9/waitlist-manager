package model

import (
	"backend-go/constant"
	"backend-go/vk"
	"context"
)

const KeyChairList = "chair-list"

type PlaceStatus struct {
	ChairList []string `json:"chair_list"`
	QueueList []int64  `json:"queue_list"`
	MaxChair  int64    `json:"max_chair"`
}

func GetPlaceStatus(ctx context.Context) (PlaceStatus, error) {
	result := PlaceStatus{ChairList: []string{}}

	chairList, err := GetChairStatus(ctx)
	if err != nil {
		return result, err
	}

	queueList, _ := GetQueueList(ctx)

	result.ChairList = chairList
	result.QueueList = queueList
	result.MaxChair = constant.MAX_CHAIR

	return result, nil
}

func GetChairStatus(ctx context.Context) ([]string, error) {
	cmd := vk.B().Lrange().Key(KeyChairList).Start(0).Stop(-1).Build()
	resp := vk.Client().Do(ctx, cmd)

	if resp.Error() != nil {
		return nil, resp.Error()
	}

	return resp.AsStrSlice()
}

func ChairListPop(ctx context.Context) (string, error) {
	cmd := vk.B().Lpop().Key(KeyChairList).Count(1).Build()
	resp := vk.Client().Do(ctx, cmd)

	if resp.Error() != nil {
		return "", resp.Error()
	}

	arr, err := resp.AsStrSlice()
	return arr[0], err
}
