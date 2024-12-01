package main

import (
	"backend-go/model"
	"context"
	"fmt"
	"time"
)

func RunQueueServerWorker() {
	ctx, cancel := context.WithCancel(context.Background())
	fmt.Println("starting queue worker")
	defer fmt.Println("stopping queue worker")
	defer cancel()

	for {
		time.Sleep(1 * time.Second)
		placeStatus, err := model.GetPlaceStatus(ctx)
		if err != nil {
			continue
		}
		if len(placeStatus.QueueList) == 0 {
			continue
		}
		fmt.Printf("place status: %+v\n", placeStatus)

		queueId := fmt.Sprint(placeStatus.QueueList[0])
		headQueueInfo, err := model.GetQueueInfo(ctx, queueId)
		if err != nil {
			continue
		}
		fmt.Printf("head queue %s info: %+v\n", queueId, headQueueInfo)

		readyQueueId, _ := model.GetReadyQueue(ctx)
		if queueId == readyQueueId {
			fmt.Printf("queue %s ready alr, waiting for checkin\n", queueId)
			continue
		}

		chairAvailable := MAX_CHAIR_LEN - len(placeStatus.ChairList)
		if chairAvailable < int(headQueueInfo.Number) {
			fmt.Printf("queue %s needs %d but only %d available, waiting\n", queueId, headQueueInfo.Number, chairAvailable)
			continue
		}
		fmt.Printf("setting queue %s to ready\n", queueId)
		err = model.SetReadyQueue(ctx, queueId)
		if err != nil {
			fmt.Printf("error setting queue %s ready: %v\n", queueId, err)
			continue
		}
	}

}

func RunChairServerWorker() {
	ctx, cancel := context.WithCancel(context.Background())
	fmt.Println("starting chair worker")
	defer fmt.Println("stopping chair worker")
	defer cancel()

	for {
		time.Sleep(3 * time.Second)
		chairList, err := model.GetChairStatus(ctx)
		if err != nil {
			continue
		}
		if len(chairList) == 0 {
			// fmt.Println("empty chair, waiting")
			continue
		}

		chair, err := model.ChairListPop(ctx)
		if err != nil {
			fmt.Printf("failed popping chair: %v", err)
			continue
		}
		fmt.Printf("serve chair %v\n", chair)
	}
}
