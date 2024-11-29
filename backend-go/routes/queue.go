package routes

import (
	"backend-go/model"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
)

type PostQueueBody struct {
	PartyName   string `json:"party_name"`
	PartyNumber int    `json:"party_number"`
}

type PostQueueResp struct {
	QueueNumber int64   `json:"queue_number"`
	QueueList   []int64 `json:"queue_list"`
}

func PostQueue(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	bodyJson, err := io.ReadAll(r.Body)
	if err != nil {
		http.Error(w, fmt.Sprint(err), http.StatusBadRequest)
		return
	}

	var body PostQueueBody
	json.Unmarshal(bodyJson, &body)

	queueNumber, err := model.AddNewQueue(ctx, body.PartyName, body.PartyNumber)
	if err != nil {
		http.Error(w, fmt.Sprint(err), http.StatusInternalServerError)
		return
	}
	queueList, _ := model.GetQueueList(ctx)

	resp := PostQueueResp{QueueNumber: queueNumber, QueueList: queueList}
	w.WriteHeader(201)
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(resp)
}
