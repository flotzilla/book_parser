package utils

import (
	"book_parser/src"
	"encoding/json"
	"os"
)

func GetConfig(configFile string) (*src.Config, error) {
	file, err := os.Open(configFile)

	defer func() {
		if file != nil {
			file.Close()
		}
	}()
	if err != nil {
		return nil, err
	}

	dec := json.NewDecoder(file)
	conf := src.Config{}

	err = dec.Decode(&conf)
	if err != nil {
		return nil, err
	}

	return &conf, nil
}

func IsStringInSlice(a string, list []string) bool {
	for _, b := range list {
		if b == a {
			return true
		}
	}
	return false
}
