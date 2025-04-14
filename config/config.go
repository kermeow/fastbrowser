package config

import (
	"os"

	"github.com/BurntSushi/toml"
)

type Config struct {
	GameDir    string   `toml:"game_dir"`
	SearchDirs []string `toml:"search_dirs"`
}

func Default() *Config {
	return &Config{
		GameDir:    "C:\\Program Files (x86)\\FastGH3\\",
		SearchDirs: make([]string, 0),
	}
}

func Load() (*Config, error) {
	d := Default()

	if _, err := os.Stat("fastbrowser.toml"); os.IsNotExist(err) {
		return d, nil
	}

	_, err := toml.DecodeFile("fastbrowser.toml", &d)
	return d, err
}
