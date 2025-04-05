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
	storage "github.com/rgglez/storage/storage"

	config "github.com/rgglez/go-config"
)

type Configuracion struct {
	ID       string `yaml:"ID"`
	Referrer string `yaml:"REFERRER"`
}

func main() {
	// Creating the storage service
	cnn := "oss://test/?credential=hmac:Secret123:Secret123&endpoint=http://127.0.0.1:9090&name=test"
	s := storage.NewStorage(cnn)

	// Setting up the configuration loader
	c := config.NewConfigurator(&config.Config{
		Referrer: "",
		Stage:    "",
		File:     "config.yaml",
	}, s)

	// Loading YAML file into a map...
	var cfgMap map[string]interface{}
	err := c.Load(&cfgMap)
	if err != nil {
		log.Fatal(err)
	}

	fmt.Println("Loading YAML file into a map...")
	pretty.Println(cfgMap)

	// Loading YAML file into a struct...
	var cfgStruct Configuracion
	err = c.Load(&cfgStruct)
	if err != nil {
		log.Fatal(err)
	}

	// Here you could validate the struct, for example...

	fmt.Println("Loading YAML file into a struct...")
	pretty.Println(cfgStruct)
}
