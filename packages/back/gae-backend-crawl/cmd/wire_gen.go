// Code generated by Wire. DO NOT EDIT.

//go:generate go run -mod=mod github.com/google/wire/cmd/wire
//go:build !wireinject
// +build !wireinject

package main

import (
	"gae-backend-crawl/bootstrap"
	"gae-backend-crawl/crawl"
	"gae-backend-crawl/executor"
	"github.com/google/wire"
)

// Injectors from wire.go:

// InitializeApp init application.
func InitializeApp() (*bootstrap.Application, error) {
	env := bootstrap.NewEnv()
	databases := bootstrap.NewDatabases(env)
	poolsFactory := bootstrap.NewPoolFactory()
	client := bootstrap.NewBloomFilter()
	kafkaConf := bootstrap.NewKafkaConf(env)
	repoCrawl := crawl.NewRepoCrawl(client, kafkaConf, env)
	crawlExecutor := executor.NewCrawlExecutor(repoCrawl)
	application := &bootstrap.Application{
		Env:           env,
		Databases:     databases,
		PoolsFactory:  poolsFactory,
		Bloom:         client,
		CrawlExecutor: crawlExecutor,
		KafkaConf:     kafkaConf,
	}
	return application, nil
}

// wire.go:

var appSet = wire.NewSet(bootstrap.NewEnv, bootstrap.NewKafkaConf, bootstrap.NewDatabases, bootstrap.NewPoolFactory, bootstrap.NewBloomFilter, crawl.NewRepoCrawl, executor.NewCrawlExecutor, wire.Struct(new(bootstrap.Application), "*"))