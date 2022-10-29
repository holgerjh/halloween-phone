package config

import (
	"fmt"
	"os"

	"github.com/holgerjh/halloween-phone/pulse"

	"gopkg.in/yaml.v2"
)

// Config holds the program's configuration
type Config struct {
	// How long to wait before playing sound after someone picked up the call in seconds
	SilenceStartOfCall int `yaml:"silenceStartOfCall"`

	// Delay between tracks in seconds
	MinWait int `yaml:"minWait"`
	MaxWait int `yaml:"maxWait"`

	// Cooldown for each track before playing it again in seconds
	TrackCooldown int `yaml:"trackCooldown"`

	// Folder to load .wav files from
	TrackFolder string `yaml:"trackFolder"`

	// Virtual pulse microphone configuration
	Mic pulse.PulseMicConfig `yaml:"mic"`
}

// LoadConfig loads configuration a yaml file pointed to by path
func LoadConfig(path string) (*Config, error) {
	raw, err := os.ReadFile(path)
	if err != nil {
		return nil, fmt.Errorf("unable to load config from %s. Make sure path exists. Error is %e", path, err)
	}
	cfg := &Config{}
	err = yaml.Unmarshal(raw, cfg)
	if err != nil {
		return nil, err
	}
	return cfg, nil
}
