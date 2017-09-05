package config

import (
	"flag"
	"io/ioutil"
	"log"
	"os"

	"github.com/caarlos0/env"
	"github.com/naoina/toml"
)

// TomlConfig - config structure
type TomlConfig struct {
	General struct {
		Listen string
	}
	Ami struct {
		Host     string `env:"AMI_HOST"`
		Port     int    `env:"AMI_PORT"`
		Username string `env:"AMI_USER"`
		Password string `env:"AMI_PASS"`
	} `toml:"ami"`
	DB struct {
		Host     string `env:"DB_HOST"`
		Database string `env:"DB_DBNAME"`
		Username string `env:"DB_USER"`
		Password string `env:"DB_PASS"`
		Debug    bool
	} `toml:"db"`
	Asterisk struct {
		Context         string `toml:"call-context" env:"ASTERISK_CONTEXT"`
		PlaybackContext string `toml:"playback-context" env:"ASTERISK_PLAYBACK_CONTEXT"`
	}
}

var config *TomlConfig
var logger = log.New(os.Stdout, "", log.Ldate|log.Ltime|log.LUTC|log.Lshortfile)

// GetConfig - get current TOML-config
func GetConfig() *TomlConfig {
	if config == nil {
		config = new(TomlConfig).readConfig()
		loadExtEnv(&config.Ami)
		loadExtEnv(&config.DB)
		loadExtEnv(&config.Asterisk)
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

func loadExtEnv(cfg interface{}) {
	err := env.Parse(cfg)
	if err != nil {
		logger.Printf("Error: in parsing ENV config (via env) %+v\n", err)
	}
}
