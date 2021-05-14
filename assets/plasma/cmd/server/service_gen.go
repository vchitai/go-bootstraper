package main

import (
	"context"
	"fmt"
	"log"
	"os"
	"os/signal"
	"sync"
	"syscall"
	"time"

	"domain_dash/service_dash/configs"
)

type masterRoutine struct {
	wg        sync.WaitGroup
	shutDowns []func(ctx context.Context)
	stopWait  sync.Once
}

func (mr *masterRoutine) fork(work func(), shutDown func(ctx context.Context)) {
	mr.wg.Add(1)
	go func() {
		defer mr.wg.Done()
		defer func() {
			if r := recover(); r != nil {
				log.Println("recovered", r)
			}
		}()
		work()
	}()

	mr.shutDowns = append(mr.shutDowns, shutDown)
}

func (mr *masterRoutine) shutdown() {
	mr.stopWait.Do(func() {
		if len(mr.shutDowns) == 0 {
			return
		}

		var completed = make(chan bool, 1)
		// Wait for maximum 15s
		ctx, cancel := context.WithDeadline(context.TODO(), time.Now().Add(15*time.Second))
		defer cancel()

		go func() {
			select {
			case <-ctx.Done():
				log.Fatal("Force shutdown due to timeout")
			case <-completed:
				log.Println("Shutdown in time")
			}
		}()

		for _, shutDown := range mr.shutDowns {
			mr.wg.Add(1)
			go func(f func(ctx context.Context), ctx context.Context) {
				defer mr.wg.Done()
				defer func() {
					if r := recover(); r != nil {
						log.Println("recovered", r)
					}
				}()
				f(ctx)
			}(shutDown, ctx)
		}

		mr.wg.Wait()
		completed <- true
	})
}

func (mr *masterRoutine) shutdownWait() {
	osSignal := make(chan os.Signal, 1)
	signal.Notify(osSignal, syscall.SIGINT, syscall.SIGTERM)
	<-osSignal
}

func loadConfig() (*configs.Config, error) {
	var cfg configs.Config
	if err := configs.Load(&cfg, configs.DefaultConfig); err != nil {
		return nil, fmt.Errorf("loading configs %w", err)
	}

	return &cfg, nil
}

func mustLoadConfig() *configs.Config {
	cfg, err := loadConfig()
	if err != nil {
		panic(err)
	}
	return cfg
}
