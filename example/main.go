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
