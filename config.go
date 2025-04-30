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
	"crypto/sha1"
	"fmt"
	"io"
	"net/url"
	"os"
	"reflect"
	"regexp"
	"strings"

	_ "github.com/rgglez/go-storage/services/oss/v3"
	"github.com/rgglez/storage/storage"

	yaml "gopkg.in/yaml.v3"
)

//-----------------------------------------------------------------------------

type Configurator struct {
	Storage    *storage.Storage
	ConfigFile string
	TmpDir     string
}

//-----------------------------------------------------------------------------

type Config struct {
	Referrer string
	Stage    string
	File     string
	TmpDir   string
}

//-----------------------------------------------------------------------------

func NewConfigurator(cfg *Config, store *storage.Storage) *Configurator {
	var domain string
	var stage string
	var file string
	var tmpDir string

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

	// The temporary dir
	if cfg.TmpDir == "" {
		tmpDir = os.TempDir()
	}

	// Remove duplicated /
	re := regexp.MustCompile(`(\/)+`)
	path = re.ReplaceAllStringFunc(path, func(m string) string {
		return "/"
	})

	return &Configurator{
		Storage:    store,
		ConfigFile: path,
		TmpDir:     tmpDir,
	}
}

//-----------------------------------------------------------------------------

func (c *Configurator) Load(config interface{}) error {
	// Local file path in the system's temporary directory
	if c.TmpDir == "" {
		c.TmpDir = os.TempDir()
	}

	h := sha1.New()
	io.WriteString(h, c.ConfigFile)
	tmpFilePath, err := os.CreateTemp(c.TmpDir, "cfg_"+fmt.Sprintf("%x", h.Sum(nil)))
	if err != nil {
		return err
	}
	defer os.Remove(tmpFilePath.Name())

	// Copy the remote YAML file to a local temporary file for parsing...
	if err := c.Storage.Read(c.ConfigFile, tmpFilePath.Name()); err != nil {
		return err
	}

	// Read the local file
	data, err := os.ReadFile(tmpFilePath.Name())
	if err != nil {
		return err
	}

	// Close (and remove) the file when done
	err = tmpFilePath.Close()
	if err != nil {
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
	err = yaml.Unmarshal(data, v.Interface())
	if err != nil {
		return err
	}

	return nil
}
