package main

import (
	"L0/internal/repository"
	"L0/internal/server"
	"L0/internal/util"

	"github.com/rs/zerolog/log"
)

func main() {
	config := util.NewConfig()
	cache := repository.NewCache()
	repository := repository.NewRepository(config, cache)

	server := server.NewServer(config, repository, cache)
	if err := server.Start(); err != nil {
		log.Err(err).Msg("the server cannot be started")
	}
}
