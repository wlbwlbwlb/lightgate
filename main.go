package main

import (
	"context"
	"os/signal"
	"syscall"
	"time"

	"github.com/wlbwlbwlb/lightgate/config"
	_ "github.com/wlbwlbwlb/lightgate/mysub"
	"github.com/wlbwlbwlb/lightgate/mytcp"

	"github.com/DarthPestilane/easytcp"
	"github.com/wlbwlbwlb/log"
	"github.com/wlbwlbwlb/mq"
)

func main() {
	//fmt.Println("hello world")

	// Create context that listens for the interrupt signal from the OS.
	ctx, stop := signal.NotifyContext(context.Background(), syscall.SIGINT, syscall.SIGTERM)
	defer stop()

	logger := log.Logger()
	defer logger.Sync()

	if e := mq.Init(mq.Lookupd(config.TOML.Nsq.Lookupd),
		mq.Nsqd(config.TOML.Nsq.Nsqd),
	); e != nil {
		log.Fatal(e.Error())
	}
	defer func() {
		mq.StopConsumers()
		mq.StopProducer()
	}()

	serve, _ := mytcp.Init(mytcp.Writer(log.Writer()))

	go func() {
		if e := serve.Run(config.TOML.Addr); e != nil && e != easytcp.ErrServerStopped {
			log.Fatalf("listen: %s\n", e)
		}
	}()

	// Listen for the interrupt signal.
	<-ctx.Done()

	// Restore default behavior on the interrupt signal and notify user of shutdown.
	stop()
	log.Info("shutting down gracefully, press Ctrl+C again to force")

	mq.StopConsumers()

	// The context is used to inform the server it has 5 seconds to finish
	// the request it is currently handling
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	if e := serve.Stop(); e != nil {
		log.Fatal("serve stop: ", e)
	}

	log.Info("serve exiting")
}
