package main

import (
	"github.com/pboehm/ddns/backend"
	"github.com/pboehm/ddns/frontend"
	"github.com/pboehm/ddns/shared"
	"golang.org/x/sync/errgroup"
	"log"
)

var serviceConfig *shared.Config = &shared.Config{}

func init() {
	serviceConfig.Initialize()
}

func main() {
	serviceConfig.Validate()

	redis := shared.NewRedisBackend(serviceConfig)
	defer redis.Close()

	var group errgroup.Group

	group.Go(func() error {
		lookup := backend.NewHostLookup(serviceConfig, redis)
		return backend.NewBackend(serviceConfig, lookup).Run()
	})

	group.Go(func() error {
		return frontend.NewFrontend(serviceConfig, redis).Run()
	})

	if err := group.Wait(); err != nil {
		log.Fatal(err)
	}
}
