package main

import (
	"context"
	"fmt"
	"io"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/wl955/lightgate/config"
	_ "github.com/wl955/lightgate/mq"
	"github.com/wl955/lightgate/tcpx"

	"github.com/DarthPestilane/easytcp"
	"github.com/wl955/log"
	"github.com/wl955/nsqx"
)

func main() {
	//fmt.Println("hello world")

	// Create context that listens for the interrupt signal from the OS.
	ctx, stop := signal.NotifyContext(context.Background(), syscall.SIGINT, syscall.SIGTERM)
	defer stop()

	defer func() {
		log.Sync()
	}()

	w := io.MultiWriter(os.Stdout, log.Writer())

	e := nsqx.Init(nsqx.Lookupd(config.TOML.Nsq.Lookupd),
		nsqx.Nsqd(config.TOML.Nsq.Nsqd),
	)
	if e != nil {
		log.Fatal(e.Error())
	}
	defer func() {
		nsqx.Stop()
	}()

	serve, _ := tcpx.Init(tcpx.Writer(w))

	go func() {
		if err := serve.Run(fmt.Sprintf(":%d", config.TOML.Port)); err != nil && err != easytcp.ErrServerStopped {
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
