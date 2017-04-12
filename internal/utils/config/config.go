package config

import (
	"flag"
	"io/ioutil"
	"os"

	"github.com/naoina/toml"
	"log"
)

// TomlConfig - config structure
type TomlConfig struct {
	General struct {
		Listen string
	}
	Ami struct {
		Host     string
		Port     int
		Username string
		Password string
	} `toml:"ami"`
	Asterisk struct {
		Context string `toml:"call-context"`
		PlaybackContext string `toml:"playback-context"`
	}
}

var config *TomlConfig
var logger = log.New(os.Stdout, "", log.Ldate|log.Ltime|log.LUTC|log.Lshortfile)

// GetConfig - get current TOML-config
func GetConfig() *TomlConfig {
	if config == nil {
		config = new(TomlConfig).readConfig()
	}

	return config
}

func (t *TomlConfig) readConfig() *TomlConfig {
	confFlag := t.readFlags()
	f, err := os.Open(*confFlag)
	if err != nil {
		logger.Printf("Can't open file: %v", err)
	}

	defer f.Close()

	buf, err := ioutil.ReadAll(f)
	if err != nil {
		logger.Printf("IO read error: %v", err)
	}

	if err := toml.Unmarshal(buf, t); err != nil {
		logger.Printf("TOML error: %v", err)
	}

	return t
}

func (t *TomlConfig) readFlags() (confFlag *string) {
	confFlag = flag.String("conf", "api.conf", "api.conf")
	flag.Parse()
	return
}
