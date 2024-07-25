package main

import (
	"net/http"

	"github.com/sudeeya/metrics-harvester/internal/repository/storage"
	"github.com/sudeeya/metrics-harvester/internal/router"
)

func main() {
	memStorage := storage.NewMemStorage()
	router := router.NewRouter(memStorage)
	err := http.ListenAndServe(":8080", router)
	if err != nil {
		panic(err)
	}
}
