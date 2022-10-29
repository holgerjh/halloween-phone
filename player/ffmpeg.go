package player

import (
	"fmt"
	"log"
	"os"
	"os/exec"

	"github.com/holgerjh/halloween-phone/pulse"
	"github.com/holgerjh/halloween-phone/statemachine"
)

const ffmpegBinary = "ffmpeg"

// PipeThroughFFMPEG uses ffmpeg to pipe file into micFile
// Playing is started as soon as the startPlaying channel gets closed
// Closing the terminate channel will stop playing and temrinate ffmpeg
// cfg must contain information about the microphone such as encoding
func PipeThroughFFMPEG(micFile string, file string, cfg *pulse.PulseMicConfig, startPlaying, terminate <-chan int) (<-chan int, error) {
	log.Printf("Looking for %s", ffmpegBinary)

	ffmpegFullPath, err := exec.LookPath(ffmpegBinary)
	if err != nil {
		return nil, err
	}
	log.Printf("Opening virtual mic file %s", micFile)
	output, err := os.OpenFile(micFile, os.O_APPEND|os.O_WRONLY, 0666)
	if err != nil {
		return nil, err
	}

	args := assembleFFMPEGArgs(micFile, file, cfg)
	log.Printf("Running %s with arguments %v into %s", ffmpegFullPath, args, micFile)
	cmd := exec.Command(ffmpegFullPath, args...)
	cmd.Stdout = output
	cmd.Stderr = os.Stderr //output

	terminated := make(chan int)
	finishedPlaying := make(chan int)

	sm := &statemachine.Statemachine{Name: "ffmpeg"}
	nodeReady := &statemachine.Node{
		Name: "ffmpeg ready",
		OnEnter: func() *statemachine.Node {
			log.Printf("ffmpeg: Ready to play, waiting for signal")
			return nil
		},
	}
	sm.AddNode(nodeReady)
	nodePlaying := &statemachine.Node{
		Name: "ffmpeg playing",
		OnEnter: func() *statemachine.Node {
			log.Printf("ffmpeg: Playing audio (async)")
			go func() {
				err = cmd.Run()
				if err != nil {
					log.Printf("ERROR: Failed playing track %s with error %e", file, err)
				}
				close(finishedPlaying)
			}()
			return nil
		},
	}
	sm.AddNode(nodePlaying)
	nodeExit := &statemachine.Node{
		Name: "ffmpeg finished",
		OnEnter: func() *statemachine.Node {
			return nil
		},
	}
	sm.AddNode(nodeExit)
	nodeShutdown := &statemachine.Node{
		Name: "ffmpeg shutdown",
		OnEnter: func() *statemachine.Node {
			log.Printf("stopping ffmpeg")
			tryKillProcess(cmd) // TODO: race condition: process might be nil (and if we ignore it it might still run afterwards)
			output.Close()
			close(terminated)
			return nodeExit
		},
	}
	sm.AddNode(nodeShutdown)

	sm.Start = nodeReady
	sm.Exit = nodeExit

	sm.AddTransition(&statemachine.Transition{
		From:      nodeReady,
		To:        nodePlaying,
		Condition: startPlaying,
	})

	sm.AddTransition(&statemachine.Transition{
		From:      nodeReady,
		To:        nodeShutdown,
		Condition: terminate,
	})

	sm.AddTransition(&statemachine.Transition{
		From:      nodePlaying,
		To:        nodeShutdown,
		Condition: finishedPlaying,
	})

	sm.AddTransition(&statemachine.Transition{
		From:      nodePlaying,
		To:        nodeShutdown,
		Condition: terminate,
	})

	go sm.Run()

	return terminated, nil

}

func assembleFFMPEGArgs(micFile, file string, cfg *pulse.PulseMicConfig) []string {
	return []string{"-re", "-i", file, "-f", cfg.Format, "-ar", fmt.Sprintf("%d", cfg.Rate), "-ac", fmt.Sprintf("%d", cfg.Channels), "-"}
}
