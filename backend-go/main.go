package main

import (
	"backend-go/model"
	"backend-go/routes"
	"context"
	"fmt"
	"net/http"
	"os"
	"sync"
	"time"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/go-chi/cors"
)

func main() {
	r := chi.NewRouter()
	r.Use(middleware.Logger)
	r.Use(cors.AllowAll().Handler)
	r.Use(middleware.RequestID)

	r.Route("/api", func(r chi.Router) {
		r.Get("/", func(w http.ResponseWriter, r *http.Request) {
			fmt.Println(r.Context().Value(middleware.RequestIDKey))
			w.Write([]byte("welcome"))
		})

		r.Get("/stream-place-status", routes.StreamPlaceStatus)
		r.Post("/queue", routes.PostQueue)
		r.Get("/queue/{queueId}", routes.GetQueue)

		r.Get("/queue/{queueId}/stream-status", routes.StreamQueueStatus)
	})

	PORT := os.Getenv("BE_PORT")

	var wg sync.WaitGroup
	wg.Add(1)
	go func() {
		http.ListenAndServe(fmt.Sprintf("0.0.0.0:%s", PORT), r)
		wg.Done()
	}()

	wg.Add(1)
	go func() {
		RunQueueWorker()
		wg.Done()
	}()

	fmt.Printf("running server on :%s\n", PORT)

	wg.Wait()
}

const MAX_CHAIR_LEN = 10

func RunQueueWorker() {
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
		fmt.Printf("head queue info: %+v\n", headQueueInfo)

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
