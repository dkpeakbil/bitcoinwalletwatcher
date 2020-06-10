package main

import (
	"context"
	"flag"
	"log"
	"os"
	"os/signal"
	"syscall"

	w "github.com/dkpeakbil/bitcoinwalletwatcher"
)

func main() {
	var (
		config = flag.String("config", "config_default.json", "default configuration file")
	)
	flag.Parse()

	cfg, err := w.NewConfig(*config)
	if err != nil {
		panic(err)
	}

	w, err := w.NewWatcher(cfg)
	if err != nil {
		panic(err)
	}
	w.SetCallback(func(addr string, amount int) {
		log.Printf("%s got %d satoshi.\n", addr, amount)
	})

	go w.Run(context.Background())

	c := make(chan os.Signal, 1)
	signal.Notify(c, os.Interrupt, syscall.SIGABRT, syscall.SIGINT, syscall.SIGQUIT)

	select {
	case s := <-c:
		log.Printf("got signal: %v\n", s)
		w.Stop()
		break
	}

	log.Print("Watcher has stopped")
}
