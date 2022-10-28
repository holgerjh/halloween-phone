package main

import (
	"github.com/holgerjh/halloween-phone/config"
	"github.com/holgerjh/halloween-phone/player"
	"github.com/holgerjh/halloween-phone/pulse"
	"github.com/holgerjh/halloween-phone/tracks"
	"log"
	"math/rand"
	"os"
	"os/signal"
	"path/filepath"
	"sync"
	"time"
)

func main() {
	log.Println("Halloween Phone started")
	home, err := os.UserHomeDir()
	if err != nil {
		panic(err)
	}
	config_path := filepath.Join(home, ".config", "halloween-phone.cfg")
	cfg, err := config.LoadConfig(config_path)
	if err != nil {
		panic(err)
	}
	log.Printf("Running with config %+v", cfg)

	rand.Seed(int64(time.Now().UnixMilli()))

	log.Printf("Initializing pulseaudio")
	pulseData, err := pulse.SetupPulseDevice(cfg.Mic)
	if err != nil {
		panic(err)
	}

	log.Printf("Loading track db")
	db, err := tracks.LoadTracksIntoNewDB(cfg.TrackFolder)
	if err != nil {
		panic(err)
	}
	log.Printf("Loaded track db")

	shutdown := make(chan int)
	wg := sync.WaitGroup{}

	startPlayThread(pulseData.VirtualMicPath, &wg, db, cfg, shutdown)

	blockUntilSigINTandSendTerminate(shutdown, &wg)
	log.Printf("Cleaning up pulse audio")
	err = pulse.UnloadPulse(pulseData, cfg.Mic)
	if err != nil {
		panic(err)
	}
	log.Printf("All done, goodbye!")

}

func blockUntilSigINTandSendTerminate(shutdown chan int, wg *sync.WaitGroup) {
	c := make(chan os.Signal, 1)
	signal.Notify(c, os.Interrupt)
	for range c {
		log.Printf("Caught SIGINT, terminating")
		close(shutdown)
		break
	}
	log.Printf("Waiting for child threads to terminate")
	wg.Wait()

}

func startPlayThread(micFile string, wg *sync.WaitGroup, db *tracks.TrackDB, cfg *config.Config, shutdown <-chan int) {
	log.Printf("Starting player thread")
	wg.Add(1)
	go func() {
		player.Loop(micFile, db, cfg, shutdown)
		wg.Done()
	}()
}
