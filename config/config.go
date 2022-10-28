package config

import (
	"fmt"
	"github.com/holgerjh/halloween-phone/pulse"
	"os"

	"gopkg.in/yaml.v2"
)

// Config foo
type Config struct {
	SilenceStartOfCall int                  `yaml:"silenceStartOfCall"`
	MinWait            int                  `yaml:"minWait"`
	MaxWait            int                  `yaml:"maxWait"`
	TrackCooldown      int                  `yaml:"trackCooldown"`
	TrackFolder        string               `yaml:"trackFolder"`
	Mic                pulse.PulseMicConfig `yaml:"mic"`
}

func stubConfig() *Config {
	return &Config{
		SilenceStartOfCall: 2,
		MinWait:            10,
		MaxWait:            40,
		TrackCooldown:      60,
		TrackFolder:        "/opt/halloween-phone/wav",
		Mic: pulse.PulseMicConfig{
			Format:   "s16le",
			Rate:     8000,
			Channels: 1,
			Dir:      "",
		},
	}
}

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
