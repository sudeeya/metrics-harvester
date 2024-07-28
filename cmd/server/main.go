package main

import (
	"net/http"

	"github.com/sudeeya/metrics-harvester/internal/handlers"
	"github.com/sudeeya/metrics-harvester/internal/repository/storage"
	"github.com/sudeeya/metrics-harvester/internal/router"
)

func main() {
	memStorage := storage.NewMemStorage()
	router := router.NewRouter(memStorage)
	router.HandleFunc("/update/gauge/", handlers.CreateGaugeHandler(memStorage))
	router.HandleFunc("/update/counter/", handlers.CreateCounterHandler(memStorage))
	router.HandleFunc("/update/", handlers.BadRequest)
	router.HandleFunc("/", http.NotFound)
	err := http.ListenAndServe(":8080", router)
	if err != nil {
		panic(err)
	}
}
