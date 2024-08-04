package main

import (
	"github.com/sudeeya/metrics-harvester/internal/repository/storage"
	"github.com/sudeeya/metrics-harvester/internal/server"
)

func main() {
	cfg, err := server.NewConfig()
	if err != nil {
		panic(err)
	}
	memStorage := storage.NewMemStorage()
	server := server.NewServer(cfg, memStorage)
	server.Run()
}
