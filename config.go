package config

import (
	"errors"
	"fmt"
	"path/filepath"

	"log"
	"net/url"
	"os"
	"regexp"
	"strings"

	_ "github.com/rgglez/go-storage/services/oss/v3"
	"github.com/rgglez/storage/storage"

	yaml "gopkg.in/yaml.v3"
)

func LoadConfig(args map[string]string) (map[string]interface{}, error) {
	var config map[string]interface{}

	cnn := ""
	domain := ""
	stage := "prod"
	file := "cfg.yaml"
	path := "/tmp/"

	for key, value := range args {
		switch key {
		case "cnn":
			cnn = value
		case "stage":
			stage = value
		case "file":
			file = value
		case "path":
			path = value
		case "referrer":
			url, err := url.Parse(value)
			if err != nil {
				panic(err)
			}
			hostname := url.Hostname()
			parts := strings.Split(hostname, ":")
			domain = parts[0]
			semiPath := "/tmp/" + domain
			if _, err := os.Stat(semiPath); errors.Is(err, os.ErrNotExist) {
				if err != nil {
					err := os.Mkdir(semiPath, os.ModePerm)
					if err != nil {
						panic(err)
					}
				}
			}
		}
	}

	if strings.HasPrefix(cnn, "file://") {
		thePath := strings.TrimPrefix(cnn, "file://")
		path = filepath.Dir(thePath)
	}

	path = path + "/" + domain + "/" + file
	re := regexp.MustCompile(`(\/)+`)
	path = re.ReplaceAllStringFunc(path, func(m string) string {
		return "/"
	})

	var data []byte

	fileInfo, err := os.Stat(path)
	if (err != nil || fileInfo.Size() == 0) && !strings.HasPrefix(cnn, "file://") {
		// no existe el archivo local y es un recurso remoto, descargar
		fmt.Println("*** AQUI ***")
		fmt.Println(cnn)
		fmt.Println("*** AQUI ***")
		store := storage.NewStorage(cnn)

		remotePath := domain + "/" + stage + "/" + file
		re := regexp.MustCompile(`(\/)+`)
		remotePath = re.ReplaceAllStringFunc(remotePath, func(m string) string {
			return "/"
		})
		remotePath = strings.TrimLeft(remotePath, "/")

		err := store.Read(remotePath, path)
		if err != nil {
			panic(remotePath)
		}
	}

	fileInfo, err = os.Stat(path)
	if err == nil {
		mode := fileInfo.Mode()
		if mode.IsRegular() {
			data, err = os.ReadFile(path)
			if err != nil {
				panic(err)
			}
		}
	} else {
		os.Remove(path)
		log.Fatalf("Cannot open config file or file not downloaded correctly: %s", path)
	}

	err = yaml.Unmarshal(data, &config)
	if err != nil {
		panic(err)
	}

	return config, nil
}
