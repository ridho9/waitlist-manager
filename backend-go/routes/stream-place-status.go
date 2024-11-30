package routes

import (
	"backend-go/model"
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
			placeStatus, err := model.GetPlaceStatus(ctx)
			if err == nil {
				jsonData, _ := json.Marshal(placeStatus)
				fmt.Fprintf(w, "data: %s\n\n", jsonData)
			}
			w.(http.Flusher).Flush()
		}
		time.Sleep(2 * time.Second)
	}

}
