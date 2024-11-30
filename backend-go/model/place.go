package model

import (
	"backend-go/vk"
	"context"
)

type PlaceStatus struct {
	ChairList []string `json:"chair_list"`
	QueueList []int64  `json:"queue_list"`
}

func GetPlaceStatus(ctx context.Context) (PlaceStatus, error) {
	result := PlaceStatus{ChairList: []string{}}

	chairList, err := fetchChairStatus(ctx)
	if err != nil {
		return result, err
	}

	queueList, _ := GetQueueList(ctx)

	result.ChairList = chairList
	result.QueueList = queueList

	return result, nil
}

func fetchChairStatus(ctx context.Context) ([]string, error) {
	cmd := vk.B().Lrange().Key("chair").Start(0).Stop(-1).Build()
	resp := vk.Client().Do(ctx, cmd)

	if resp.Error() != nil {
		return nil, resp.Error()
	}

	return resp.AsStrSlice()
}
