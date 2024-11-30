package main

import (
	"backend-go/routes"
	"fmt"
	"net/http"
	"os"
	"sync"

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
	})

	PORT := os.Getenv("BE_PORT")

	var wg sync.WaitGroup
	wg.Add(1)
	go func() {
		http.ListenAndServe(fmt.Sprintf("0.0.0.0:%s", PORT), r)
		wg.Done()
	}()

	fmt.Printf("running server on :%s\n", PORT)

	wg.Wait()
}
