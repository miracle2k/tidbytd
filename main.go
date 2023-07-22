package main

import (
	"fmt"
	"io/ioutil"
	"log"
	"net/http"

	"gopkg.in/yaml.v2"
)

type Device struct {
	Id       string `yaml:"id"`
	ApiToken string `yaml:"api_token"`
}

type Content struct {
	URL  string `yaml:"url"`
	Name string `yaml:"name"`
}

type Config struct {
	Devices map[string]Device `yaml:"devices"`
	Content []Content         `yaml:"content"`
}

func status(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "OK\n")
}

func main() {
	// Load config on boot
	yamlFile, err := ioutil.ReadFile("config.yaml")
	if err != nil {
		log.Fatalf("Unable to read config file: %v", err)
	}

	var cfg Config
	err = yaml.Unmarshal(yamlFile, &cfg)
	if err != nil {
		log.Fatalf("Unable to unmarshal config: %v", err)
	}

	fmt.Println("config loaded")

	fmt.Println("setup scheduler")
	schedule(cfg)

	fmt.Println("serving")
	http.HandleFunc("/status", status)
	http.ListenAndServe(":8080", nil)
}
