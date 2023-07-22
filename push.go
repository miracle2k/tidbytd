package main

import (
	"bytes"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
)

type PushOpts struct {
	DeviceId       string
	InstallationId string
	APIToken       string
	Background     bool
}

type TidbytPushJSON struct {
	DeviceID       string `json:"deviceID"`
	Image          string `json:"image"`
	InstallationID string `json:"installationID"`
	Background     bool   `json:"background"`
}

const (
	TidbytAPIPush = "https://api.tidbyt.com/v0/devices/%s/push"
)

func push(imageData []byte, opts PushOpts) error {
	payload, err := json.Marshal(
		TidbytPushJSON{
			DeviceID:       opts.DeviceId,
			Image:          base64.StdEncoding.EncodeToString(imageData),
			InstallationID: opts.InstallationId,
			Background:     opts.Background,
		},
	)
	if err != nil {
		return fmt.Errorf("failed to marshal json: %w", err)
	}

	client := &http.Client{}
	req, err := http.NewRequest(
		"POST",
		fmt.Sprintf(TidbytAPIPush, opts.DeviceId),
		bytes.NewReader(payload),
	)
	if err != nil {
		return fmt.Errorf("creating POST request: %w", err)
	}

	req.Header.Add("Authorization", fmt.Sprintf("Bearer %s", opts.APIToken))

	resp, err := client.Do(req)
	if err != nil {
		return fmt.Errorf("pushing to API: %w", err)
	}

	if resp.StatusCode != 200 {
		fmt.Printf("Tidbyt API returned status %s\n", resp.Status)
		body, _ := ioutil.ReadAll(resp.Body)
		fmt.Println(string(body))
		return fmt.Errorf("Tidbyt API returned status: %s", resp.Status)
	}

	return nil
}
