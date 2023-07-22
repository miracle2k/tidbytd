package main

import (
	"fmt"
	"github.com/go-co-op/gocron"
	"io/ioutil"
	"strings"
	"tidbyt.dev/pixlet/runtime"
	"time"
)

func schedule(config Config) {
	s := gocron.NewScheduler(time.UTC)

	cache := runtime.NewInMemoryCache()

	for _, contentItem := range config.Content {

		// Avoid loop closure problem
		item := contentItem
		_, err := s.Every(60).Seconds().Do(func() {
			fmt.Printf("process: %s\n", item.URL)

			// Detect if item.URL is a local path
			var content []byte
			var err error
			if strings.HasPrefix(item.URL, "/") {
				content, err = ioutil.ReadFile(item.URL)
				if err != nil {
					fmt.Printf("error: read: %s: %s\n", item.URL, err)
					return
				}
			} else {
				content, err = fetchUrl(item.URL)
				if err != nil {
					fmt.Printf("error: fetch: %s: %s\n", item.URL, err)
					return
				}
			}

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
			image, result := render(content, opts, renderVars)

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

			fmt.Printf("done: %s\n", item.URL)
		})

		if err != nil {
			fmt.Printf("error: job init: %s\n", err)
		}
	}

	s.StartAsync()
}
