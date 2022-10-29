package pulse

import (
	"fmt"
	"log"
	"os"
	"os/user"
	"path/filepath"

	"github.com/lawl/pulseaudio"
)

type PulseMicConfig struct {
	Address  string `yaml:"address"`
	Dir      string `yaml:"dir"`
	Format   string `yaml:"format"`
	Rate     int    `yaml:"rate"`
	Channels int    `yaml:"channels"`
	Filename string `yaml:"filename"`
}

func (p PulseMicConfig) assembleModuleArguments(file string) string {
	return fmt.Sprintf("source_name=virtmic file=%s format=%s rate=%d channels=%d",
		file, p.Format, p.Rate, p.Channels,
	)
}

func (p PulseMicConfig) connect() (*pulseaudio.Client, error) {
	log.Printf("Registering pulse client")
	client, err := pulseaudio.NewClient(p.Address)
	if err != nil {
		return nil, err
	}

	alive := client.Connected()
	if !alive {
		return nil, fmt.Errorf("not connected to pulse audio")

	}
	log.Printf("-> Connected")
	return client, nil
}

func (p *PulseMicConfig) fillDefaults() {
	p.Filename = orDefault(p.Filename, defaultMicFilename)
	p.Address = orDefault(p.Address, defaultPulseAddress)
}

const defaultMicFilename = "virtual-mic"

var defaultPulseAddress string

func init() {
	currentUser, err := user.Current()
	if err != nil {
		panic(err)
	}

	defaultPulseAddress = fmt.Sprintf("/run/user/%s/pulse/native", currentUser.Uid)
}

const pulseModule = "module-pipe-source"

func orDefault(s, fallback string) string {
	if s == "" {
		return fallback
	}
	return s
}

type PulseData struct {
	PulseModuleIDX uint32
	VirtualMicPath string
}

// SetupPulseDevice creates and registers a virtual microphone
func SetupPulseDevice(config PulseMicConfig) (*PulseData, error) {
	config.fillDefaults()

	client, err := config.connect()
	if err != nil {
		return nil, err
	}

	log.Printf("Allocating temporary directory for pulse mic file")
	dir, err := os.MkdirTemp(config.Dir, "halloweenphone")
	if err != nil {
		return nil, err
	}

	file := filepath.Join(dir, config.Filename)

	log.Printf("Loading virtual microphone pipe module")
	args := config.assembleModuleArguments(file)
	log.Printf("-> loading %s %s", pulseModule, args)
	idx, err := client.LoadModule(pulseModule, args)
	if err != nil {
		return nil, err
	}
	log.Printf("-> Success! Index: %d", idx)
	log.Printf("Closing pulse client")
	client.Close()
	return &PulseData{
		PulseModuleIDX: idx,
		VirtualMicPath: file,
	}, nil
}

// UnloadPulse unloads and deletes an existing virtual microphone
func UnloadPulse(data *PulseData, config PulseMicConfig) error {
	config.fillDefaults()
	client, err := config.connect()
	if err != nil {
		return err
	}
	log.Printf("Unloading pulse module %d", data.PulseModuleIDX)
	err = client.UnloadModule(data.PulseModuleIDX)
	if err != nil {
		return err
	}
	log.Printf("-> success!")
	log.Printf("Unlinking virtual microphone file %s", data.VirtualMicPath)
	err = os.Remove(data.VirtualMicPath)
	if err != nil {
		return err
	}
	log.Printf("-> success!")
	return nil
}
