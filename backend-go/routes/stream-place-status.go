package routes

import (
	"backend-go/model"
	"backend-go/valkey"
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"time"
)

func StreamPlaceStatus(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	w.Header().Set("Access-Control-Allow-Origin", "*")
	w.Header().Set("Access-Control-Expose-Headers", "Content-Type")

	w.Header().Set("Content-Type", "text/event-stream")
	w.Header().Set("Cache-Control", "no-cache")
	w.Header().Set("Connection", "keep-alive")

msgLoop:
	for {
		select {
		case <-ctx.Done():
			fmt.Println("connection closed!")
			break msgLoop
		default:
			placeStatus, err := getPlaceStatus(ctx)
			if err == nil {
				jsonData, _ := json.Marshal(placeStatus)
				fmt.Fprintf(w, "data: %s\n\n", jsonData)
			}
			time.Sleep(2 * time.Second)
			w.(http.Flusher).Flush()
		}
	}

}

type PlaceStatus struct {
	ChairList []string `json:"chair_list"`
	QueueList []int64  `json:"queue_list"`
}

func getPlaceStatus(ctx context.Context) (PlaceStatus, error) {
	result := PlaceStatus{ChairList: []string{}}

	chairList, err := fetchChairStatus(ctx)
	if err != nil {
		return result, err
	}

	queueList, _ := model.GetQueueList(ctx)

	result.ChairList = chairList
	result.QueueList = queueList

	return result, nil
}

func fetchChairStatus(ctx context.Context) ([]string, error) {
	cmd := valkey.B().Lrange().Key("chair").Start(0).Stop(-1).Build()
	resp := valkey.Client().Do(ctx, cmd)

	if resp.Error() != nil {
		return nil, resp.Error()
	}

	return resp.AsStrSlice()
}
