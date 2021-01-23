package config

import (
	"crypto/sha256"
	"encoding/binary"
	"encoding/hex"
	"encoding/json"
	"github.com/sirupsen/logrus"
	"os"
)

type Config struct {
	ScanExt         []string
	SkippedExt      []string
	WithCoverImages bool
}

func GetConfig(configFile string) (*Config, error) {
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
	conf := Config{}

	err = dec.Decode(&conf)
	if err != nil {
		return nil, err
	}

	return &conf, nil
}

func (cnf *Config) GetConfigHash() string {
	hasher := sha256.New()

	for _, el := range cnf.ScanExt {
		hasher.Write([]byte(el))
	}

	for _, el := range cnf.SkippedExt {
		hasher.Write([]byte(el))
	}

	var withCover byte
	if cnf.WithCoverImages {
		withCover = 1
	}

	b1 := make([]byte, 8)
	binary.PutVarint(b1, int64(withCover))
	hasher.Write(b1)

	hash := hex.EncodeToString(hasher.Sum(nil))
	logrus.Trace("Configuration hash: ", hash)
	return hash
}

func (cnf *Config) ShowConfig() {
	logrus.Debug("\tConfiguration:")
	logrus.Debug("* ScanExt: ", cnf.ScanExt)
	logrus.Debug("* SkippedExt: ", cnf.SkippedExt)
	logrus.Debug("* With Covers: ", cnf.WithCoverImages)
}
