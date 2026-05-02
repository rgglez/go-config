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

package config_test

import (
	"context"
	"errors"
	"io"
	"testing"

	config "github.com/rgglez/go-config"
	types "github.com/rgglez/go-storage/v5/types"
)

// mockStorager implements types.Storager using the required embed and a
// configurable ReadWithContext response.
type mockStorager struct {
	types.UnimplementedStorager
	content string
	err     error
}

func (m *mockStorager) ReadWithContext(_ context.Context, _ string, w io.Writer, _ ...types.Pair) (int64, error) {
	if m.err != nil {
		return 0, m.err
	}
	n, err := io.WriteString(w, m.content)
	return int64(n), err
}

// ------- Load tests -------

const sampleYAML = "ID: hello\nREFERRER: https://example.com\n"

type appConfig struct {
	ID       string `yaml:"ID"`
	Referrer string `yaml:"REFERRER"`
}

func newConfigurator(content string, err error) *config.Configurator {
	return &config.Configurator{
		Storage:    &mockStorager{content: content, err: err},
		ConfigFile: "example.com/prod/config.yaml",
	}
}

func TestLoad_IntoStruct(t *testing.T) {
	c := newConfigurator(sampleYAML, nil)
	var cfg appConfig
	if err := c.Load(&cfg); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if cfg.ID != "hello" {
		t.Errorf("ID: got %q, want %q", cfg.ID, "hello")
	}
	if cfg.Referrer != "https://example.com" {
		t.Errorf("Referrer: got %q, want %q", cfg.Referrer, "https://example.com")
	}
}

func TestLoad_IntoMap(t *testing.T) {
	c := newConfigurator(sampleYAML, nil)
	var m map[string]any
	if err := c.Load(&m); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if m["ID"] != "hello" {
		t.Errorf("ID: got %v, want %q", m["ID"], "hello")
	}
}

func TestLoad_StorageError(t *testing.T) {
	want := errors.New("storage unavailable")
	c := newConfigurator("", want)
	var cfg appConfig
	if err := c.Load(&cfg); !errors.Is(err, want) {
		t.Errorf("got %v, want %v", err, want)
	}
}

func TestLoad_InvalidYAML(t *testing.T) {
	c := newConfigurator(":\tinvalid: [yaml", nil)
	var cfg appConfig
	if err := c.Load(&cfg); err == nil {
		t.Error("expected error for invalid YAML, got nil")
	}
}

func TestLoad_NonPointer(t *testing.T) {
	c := newConfigurator(sampleYAML, nil)
	var cfg appConfig
	if err := c.Load(cfg); err == nil {
		t.Error("expected error for non-pointer, got nil")
	}
}

func TestLoad_NilPointer(t *testing.T) {
	c := newConfigurator(sampleYAML, nil)
	var cfg *appConfig
	if err := c.Load(&cfg); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if cfg == nil || cfg.ID != "hello" {
		t.Errorf("expected populated struct via nil pointer, got %v", cfg)
	}
}

// ------- NewConfigurator tests -------

func TestNewConfigurator_PathConstruction(t *testing.T) {
	store := &mockStorager{}
	c := config.NewConfigurator(&config.Config{
		Referrer: "https://example.com",
		Stage:    "prod",
		File:     "config.yaml",
	}, store)

	if c.ConfigFile != "example.com/prod/config.yaml" {
		t.Errorf("ConfigFile: got %q, want %q", c.ConfigFile, "example.com/prod/config.yaml")
	}
}

func TestNewConfigurator_DeduplicatesSlashes(t *testing.T) {
	store := &mockStorager{}
	c := config.NewConfigurator(&config.Config{
		Referrer: "https://example.com",
		Stage:    "",
		File:     "config.yaml",
	}, store)

	for i := 0; i < len(c.ConfigFile)-1; i++ {
		if c.ConfigFile[i] == '/' && c.ConfigFile[i+1] == '/' {
			t.Errorf("double slash in ConfigFile: %q", c.ConfigFile)
		}
	}
	if c.ConfigFile[0] == '/' {
		t.Errorf("leading slash in ConfigFile: %q", c.ConfigFile)
	}
}

func TestNewConfigurator_EmptyFilePanics(t *testing.T) {
	defer func() {
		if r := recover(); r == nil {
			t.Error("expected panic for empty File, got none")
		}
	}()
	config.NewConfigurator(&config.Config{
		Referrer: "https://example.com",
		Stage:    "prod",
		File:     "",
	}, &mockStorager{})
}

// ------- FindByID tests -------

type item struct {
	id   string
	name string
}

func (it item) GetID() string { return it.id }

func TestFindByID_Found(t *testing.T) {
	items := []item{{"a", "alpha"}, {"b", "beta"}, {"c", "gamma"}}
	got, err := config.FindByID(items, "b")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if got.name != "beta" {
		t.Errorf("name: got %q, want %q", got.name, "beta")
	}
}

func TestFindByID_NotFound(t *testing.T) {
	items := []item{{"a", "alpha"}}
	if _, err := config.FindByID(items, "z"); err == nil {
		t.Error("expected error for missing ID, got nil")
	}
}

func TestFindByID_EmptySlice(t *testing.T) {
	if _, err := config.FindByID([]item{}, "a"); err == nil {
		t.Error("expected error for empty slice, got nil")
	}
}
