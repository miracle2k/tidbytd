package main

import (
	"io/ioutil"
	"net/http"
	"net/url"
	"strings"

	"github.com/pkg/errors"
)

func fetchUrl(urlStr string) ([]byte, error) {
	parsedURL, err := url.Parse(urlStr)
	if err != nil {
		return nil, errors.Wrap(err, "failed to parse URL")
	}

	if parsedURL.Host == "github.com" {
		parts := strings.Split(parsedURL.Path, "/")
		if len(parts) < 5 {
			return nil, errors.New("URL path does not contain a valid file path")
		}
		urlStr = "https://raw.githubusercontent.com/" + parts[1] + "/" + parts[2] + "/" + parts[4] + "/" + strings.Join(parts[5:], "/")
	}

	resp, err := http.Get(urlStr)
	if err != nil {
		return nil, errors.Wrap(err, "failed to download file")
	}
	defer resp.Body.Close()

	if resp.StatusCode != 200 {
		return nil, errors.Errorf("failed to download file: %s", resp.Status)
	}

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, errors.Wrap(err, "failed to read file contents")
	}

	return body, nil
}
