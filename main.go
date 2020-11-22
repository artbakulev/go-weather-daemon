package main

import (
	"context"
	"flag"
	"fmt"
	"os"
	"os/signal"
	"review/openWeather"
	"strings"
	"syscall"
	"time"
)

type Config struct {
	Interval    time.Duration
	AppId       string
	Query       string
	ParsedQuery []string
}

var config = Config{}

func init() {
	flag.DurationVar(&config.Interval, "interval", time.Second*10, "")
	flag.StringVar(&config.AppId, "appId", "0a69a6d4ea256803c7c3d1bff32d5d6a", "")
	flag.StringVar(&config.Query, "query", "Moscow,London,Paris,Yekaterinburg,Dubai", "")
}

func getWeather(ctx context.Context, query string, weatherChan chan<- string) {
	ticker := time.NewTicker(config.Interval)
	defer ticker.Stop()
	for {
		select {
		case <-ctx.Done():
			return
		case <-ticker.C:
			temperature, err := api.GetWeatherForQuery(ctx, query)
			if err != nil {
				fmt.Println(err)
			}
			weatherChan <- fmt.Sprintf("Temperature from %s: %.2f", query, temperature)
		}
	}
}

func worker(ctx context.Context, signalChan <-chan os.Signal) {
	ctx, cancel := context.WithCancel(ctx)
	weatherChan := make(chan string, len(config.ParsedQuery))
	for _, query := range config.ParsedQuery {
		go getWeather(ctx, query, weatherChan)
	}
	for {
		select {
		case <-signalChan:
			cancel()
		case result := <-weatherChan:
			fmt.Println(result)
		}
	}
}

var api openWeather.API

func main() {
	flag.Parse()
	config.ParsedQuery = strings.Split(config.Query, ",")
	api = openWeather.API{
		AppId:   config.AppId,
		Timeout: time.Second,
	}
	ctx := context.Background()
	signalChan := make(chan os.Signal, 1)
	signal.Notify(signalChan, syscall.SIGTERM, syscall.SIGINT)
	worker(ctx, signalChan)
}
