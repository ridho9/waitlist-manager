package routes

import (
	"backend-go/model"
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"os"
	"time"

	"github.com/go-chi/chi/v5"
)

func StreamQueueStatus(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	w.Header().Set("Access-Control-Allow-Origin", "*")
	w.Header().Set("Access-Control-Expose-Headers", "Content-Type")

	w.Header().Set("Content-Type", "text/event-stream")
	w.Header().Set("Cache-Control", "no-cache")
	w.Header().Set("Connection", "keep-alive")

	queueId := chi.URLParam(r, "queueId")
	fmt.Println(queueId)

msgLoop:
	for {
		select {
		case <-ctx.Done():
			fmt.Println("connection closed!")
			break msgLoop
		default:
			placeStatus, _ := model.GetPlaceStatus(ctx)
			queueStatus, err := getQueueStatus(ctx, queueId)
			queueStatus.PlaceStatus = placeStatus
			if err != nil {
				fmt.Fprintf(os.Stderr, "error: %v\n", err)
			} else {
				jsonData, _ := json.Marshal(queueStatus)
				fmt.Fprintf(w, "data: %s\n\n", jsonData)
			}
		}
		w.(http.Flusher).Flush()
		time.Sleep(2 * time.Second)
	}

}

type QueueStatus struct {
	model.PlaceStatus
	Ready     bool `json:"ready"`
	CheckedIn bool `json:"checked_in"`
}

func getQueueStatus(ctx context.Context, queueId string) (QueueStatus, error) {
	queueStatus, err := model.GetQueueInfo(ctx, queueId)
	if err != nil {
		return QueueStatus{}, nil
	}
	readyQueue, err := model.GetReadyQueue(ctx)
	if err != nil {
		return QueueStatus{}, err
	}
	return QueueStatus{Ready: readyQueue == queueId, CheckedIn: queueStatus.CheckedIn}, err
}
