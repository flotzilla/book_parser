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

	conf.HasInit = true
	return &conf, nil
}
