package main

import (
	"github.com/joho/godotenv"
	"log"
	"mqtg-bot/internal"
	"os"
	"os/signal"
	"syscall"
)

func main() {
	err := godotenv.Load()
	if err != nil {
		log.Printf("Using environment variables from container environment")
	} else {
		log.Printf("Using environment variables from .env file")
	}

	bot := internal.InitTelegramBot()

	var gracefulStop = make(chan os.Signal)
	signal.Notify(gracefulStop, syscall.SIGTERM, syscall.SIGINT)

	go func() {
		sig := <-gracefulStop
		log.Printf("Caught system sig: %+v", sig)
		bot.Shutdown()
	}()

	bot.Wait()
}
