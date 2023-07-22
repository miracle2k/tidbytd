package main

import (
	"fmt"
	"github.com/go-co-op/gocron"
	"tidbyt.dev/pixlet/runtime"
	"time"
)

func schedule(config Config) {
	s := gocron.NewScheduler(time.UTC)

	cache := runtime.NewInMemoryCache()

	for _, item := range config.Content {
		_, err := s.Every(60).Seconds().Do(func() {
			fmt.Printf("render: %s\n", item.URL)

			renderVars := map[string]string{}
			opts := RenderOpts{
				Cache:         cache,
				Magnify:       1,
				RenderGif:     false,
				MaxDuration:   15000,
				SilenceOutput: false,
				Width:         64,
				Height:        32,
			}
			image, result := render(item.URL, opts, renderVars)

			if result != nil {
				fmt.Printf("error: render: %s %s\n", item.URL, result)
				return
			}

			// Now push!
			// push to each device
			for deviceId, device := range config.Devices {
				pushOpts := PushOpts{
					DeviceId:       device.Id,
					InstallationId: item.Name,
					APIToken:       device.ApiToken,
					Background:     true,
				}
				result = push(image, pushOpts)
				if result != nil {
					fmt.Printf("error: push: %s %s %s\n", item.URL, deviceId, result)
					return
				}
			}
		})

		if err != nil {
			fmt.Printf("error: job init: %s\n", err)
		}
	}

	s.StartAsync()
}
