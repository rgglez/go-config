# go-config

[![License](https://img.shields.io/badge/License-Apache_2.0-blue.svg)](https://opensource.org/licenses/Apache-2.0)
![GitHub all releases](https://img.shields.io/github/downloads/rgglez/go-config/total)
![GitHub issues](https://img.shields.io/github/issues/rgglez/go-config)
![GitHub commit activity](https://img.shields.io/github/commit-activity/y/rgglez/go-config)
[![Go Report Card](https://goreportcard.com/badge/github.com/rgglez/go-config)](https://goreportcard.com/report/github.com/rgglez/go-config)
[![GitHub release](https://img.shields.io/github/release/rgglez/go-config.svg)](https://github.com/rgglez/go-config/releases/)
![GitHub stars](https://img.shields.io/github/stars/rgglez/go-config?style=social)
![GitHub forks](https://img.shields.io/github/forks/rgglez/go-config?style=social)

`go-config` is a module which loads a YAML configuration file from a "remote" source. It supports multiple storage backends through [storage](https://github.com/rgglez/storage) and [go-storage](https://github.com/rgglez/go-storage).

The intended usage is to load the configuration from a YAML file, for a web application written in [Go](https://golang.org). The file could be local or remote (from any source supported by [go-storage](https://github.com/rgglez/go-storage)).

## Installation

```bash
go get github.com/rgglez/go-config
```

## Configuration

```go
c := config.NewConfigurator(&config.Config{
    Referrer: "",
    Stage:    "",
    File:     "config.yaml",
}, s)
```

The configuration `config.Config` struct taken as the first parameter has this properties:

* `Referrer` string, optional, the referrer as provided by the web framework you are using.
* `Stage`    string, an optional prefix for the path of the configuration. For example: "dev" or "production".
* `File`     string, required, the file name of the configuration file.
* `TmpDir`   string, optional path for a temporary directory where the remote file will be downloaded into a local temporary file.

The second parameter, `s` in the example, must be a [storage](https://github.com/rgglez/storage) object:

```go
s := storage.NewStorage(cnn)
```

Where `cnn` is a valid connection string as specified by [go-storage](https://github.com/rgglez/go-storage).

For example, for a bucket named `test` in an [ossemulator](https://github.com/aliyun/oss-emulator) server running at `localhost`, the connection string would be:

```go
cnn := "oss://test/?credential=hmac:Secret123:Secret123&endpoint=http://127.0.0.1:9090&name=test"
```

The configuration file path is formed by these components:

```go
domain + "/" + stage + "/" + file
```

where the `domain` is obtained from the `Referrer`.

Since `domain` and `stage` are optional, the file could be in the root of the remote path. For instance, if the source is a S3 bucket, the key could:

```
example.com/config.yaml
```

or

```
prod/config.yaml
```

or even

```
config.yaml
```

## Usage

The configuration file could be loaded into a struct (which should reflect the structure of your YAML file) or a map. For example:

```go
var cfgMap map[string]interface{}
err := c.Load(&cfgMap)
```

```go
var cfgStruct Configuracion
err = c.Load(&cfgStruct)
```

See the [sample](example/) code.

## Dependencies

This module uses:

* [storage](github.com/rgglez/storage)
* [go-storage](github.com/rgglez/go-storage)

and their respective dependencies.

## License

Copyright 2024 Rodolfo González González.

Released under [Apache 2.0 license](LICENSE).
