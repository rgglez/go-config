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

package config

import (
	"bytes"
	"context"
	"fmt"
	"net/url"
	"reflect"
	"regexp"
	"strings"

	"github.com/kr/pretty"
	types "github.com/rgglez/go-storage/v5/types"

	yaml "gopkg.in/yaml.v3"
)

//-----------------------------------------------------------------------------

type Configurator struct {
	Storage    types.Storager
	ConfigFile string
}

//-----------------------------------------------------------------------------

type Config struct {
	Referrer string
	Stage    string
	File     string
}

//-----------------------------------------------------------------------------

func NewConfigurator(cfg *Config, store types.Storager) *Configurator {
	var domain string
	var stage string
	var file string

	// Get the domain from the referrer
	url, err := url.Parse(cfg.Referrer)
	if err != nil {
		panic(err)
	}
	hostname := url.Hostname()
	parts := strings.Split(hostname, ":")
	domain = parts[0]

	// The stage (prod, dev, etc.)
	stage = cfg.Stage

	// The file part
	if cfg.File == "" {
		panic("the configuration file name can not be empty")
	}
	file = cfg.File

	// Construct the path
	path := domain + "/" + stage + "/" + file

	// Remove duplicated /
	re := regexp.MustCompile(`(\/)+`)
	path = re.ReplaceAllStringFunc(path, func(m string) string {
		return "/"
	})

	path = strings.TrimLeft(path, "/")

	return &Configurator{
		Storage:    store,
		ConfigFile: path,
	}
}

//-----------------------------------------------------------------------------

func (c *Configurator) Load(config interface{}) error {
	var buf bytes.Buffer
	ctx := context.Background()
	if _, err := c.Storage.ReadWithContext(ctx, c.ConfigFile, &buf); err != nil {
		pretty.Println(err)
		return err
	}

	// Get the underlying value from the interface
	v := reflect.ValueOf(config)

	// Ensure we got a pointer
	if v.Kind() != reflect.Ptr {
		return fmt.Errorf("target must be a pointer, got %T", config)
	}

	// Create a new pointer if it's nil
	if v.IsNil() {
		v.Set(reflect.New(v.Type().Elem()))
	}

	// Unmarshal into the actual pointer value
	if err := yaml.Unmarshal(buf.Bytes(), v.Interface()); err != nil {
		pretty.Println(err)
		return err
	}

	return nil
}
