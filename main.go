package main

import (
	"context"
	"os/signal"
	"syscall"
	"time"

	"github.com/wl955/lightgate/config"
	_ "github.com/wl955/lightgate/pubsub"
	"github.com/wl955/lightgate/tcpx"

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

	quit, e := mq.Init(mq.Lookupd(config.TOML.Nsq.Lookupd),
		mq.Nsqd(config.TOML.Nsq.Nsqd),
	)
	if e != nil {
		log.Fatal(e.Error())
	}
	defer quit()

	serve, _ := tcpx.Init(tcpx.Writer(log.Writer()))

	go func() {
		if err := serve.Run(config.TOML.Addr); err != nil && err != easytcp.ErrServerStopped {
			log.Fatalf("listen: %s\n", err)
		}
	}()

	// Listen for the interrupt signal.
	<-ctx.Done()

	// Restore default behavior on the interrupt signal and notify user of shutdown.
	stop()
	log.Info("shutting down gracefully, press Ctrl+C again to force")

	// The context is used to inform the server it has 5 seconds to finish
	// the request it is currently handling
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	if err := serve.Stop(); err != nil {
		log.Fatal("serve stop: ", err)
	}

	log.Info("serve exiting")
}
