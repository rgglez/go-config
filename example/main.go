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
	cnn := "oss://test/?credential=hmac:Secret123:Secret123&endpoint=http://127.0.0.1:9090&name=test"
	s := storage.NewStorage(cnn)

	c := config.NewConfigurator(&config.Config{
		Referrer: "",
		Stage:    "",
		File:     "auth.yaml",
	}, s)

	//var cfg Configuracion
	var cfg map[string]interface{}
	err := c.Load(&cfg)
	if err != nil {
		log.Fatal(err)
	}

	pretty.Println(cfg)
}
