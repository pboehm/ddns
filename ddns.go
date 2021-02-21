package main

import (
	"github.com/pboehm/ddns/backend"
	"github.com/pboehm/ddns/frontend"
	"github.com/pboehm/ddns/shared"
	"github.com/prometheus/client_golang/prometheus"
	"golang.org/x/sync/errgroup"
	"log"
)

var serviceConfig *shared.Config = &shared.Config{}

func init() {
	serviceConfig.Initialize()
}

func main() {
	serviceConfig.Validate()

	if serviceConfig.Verbose {
		log.Printf("Loaded config: %#v", serviceConfig)
	}

	redis := shared.NewRedisBackend(serviceConfig)
	defer redis.Close()

	registry := prometheus.NewRegistry()

	var group errgroup.Group

	group.Go(func() error {
		lookup := backend.NewHostLookup(serviceConfig, redis)
		return backend.NewBackend(serviceConfig, lookup, registry).Run()
	})

	group.Go(func() error {
		return frontend.NewFrontend(serviceConfig, redis, registry).Run()
	})

	if err := group.Wait(); err != nil {
		log.Fatal(err)
	}
}
