/*
Copyright 2024 Rodolfo Gonzalez Gonzalez

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

	http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/

package main

import (
	"fmt"
	"log"

	"github.com/kr/pretty"
	"github.com/rgglez/go-storage/services/oss/v3"
	"github.com/rgglez/go-storage/v5/pairs"

	config "github.com/rgglez/go-config"
)

type DBConfig struct {
	Host string `yaml:"host"`
	Port int    `yaml:"port"`
	Name string `yaml:"name"`
}

type Feature struct {
	ID      string `yaml:"id"`
	Enabled bool   `yaml:"enabled"`
}

func (f Feature) GetID() string { return f.ID }

type AppConfig struct {
	AppName  string    `yaml:"app_name"`
	Version  string    `yaml:"version"`
	Debug    bool      `yaml:"debug"`
	Database DBConfig  `yaml:"database"`
	Features []Feature `yaml:"features"`
}

func main() {
	// Connect to an OSS-compatible bucket (MinIO, Alibaba OSS, etc.).
	// Adjust credential, endpoint and name for your environment.
	_, store, err := oss.New(
		pairs.WithCredential("hmac:Secret123:Secret123"),
		pairs.WithEndpoint("http://127.0.0.1:9090"),
		pairs.WithName("test"),
	)
	if err != nil {
		log.Fatal(err)
	}

	// Path in storage will be: myapp.example.com/prod/config.yaml
	c := config.NewConfigurator(&config.Config{
		Referrer: "https://myapp.example.com",
		Stage:    "prod",
		File:     "config.yaml",
	}, store)

	// Load into a generic map.
	var cfgMap map[string]any
	if err := c.Load(&cfgMap); err != nil {
		log.Fatal(err)
	}
	fmt.Println("--- map ---")
	pretty.Println(cfgMap)

	// Load into a typed struct.
	var cfg AppConfig
	if err := c.Load(&cfg); err != nil {
		log.Fatal(err)
	}
	fmt.Println("--- struct ---")
	pretty.Println(cfg)

	// Look up a feature flag by ID.
	f, err := config.FindByID(cfg.Features, "dark-mode")
	if err != nil {
		log.Fatal(err)
	}
	fmt.Printf("dark-mode enabled: %v\n", f.Enabled)
}
