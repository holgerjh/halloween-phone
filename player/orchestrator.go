package player

import (
	"log"
	"time"

	"github.com/holgerjh/halloween-phone/pulse"
	"github.com/holgerjh/halloween-phone/statemachine"
	"github.com/holgerjh/halloween-phone/tracks"
)

const MAX_RING_TIME = 10

func playTrack(micFile string, t *tracks.Track, silenceStartOfCall int, cfg *pulse.PulseMicConfig, shutdown <-chan int) error {
	log.Printf("Will try to play %s", t.Path)
	now := time.Now()
	t.LastPlayedAt = &now

	// channels to abort calling and ffmpeg conversion
	terminateCall := make(chan int)
	terminateFFMPEG := make(chan int)

	// channel to signal FFMPEG process to start creating fake mic input
	startPlayFFMPEG := make(chan int)

	// setup async ffmpeg conversion

	ffmpegTerminated, err := PipeThroughFFMPEG(micFile, t.Path, cfg, startPlayFFMPEG, terminateFFMPEG)
	if err != nil {
		log.Printf("Unable to start ffmpeg")
		return err
	}

	// make async call
	callConnected, callTerminated, err := Call(MAX_RING_TIME, terminateCall)
	if err != nil {
		log.Printf("Unable to start call")
		// make sure to shutdown ffmpeg properly
		close(terminateFFMPEG)
		return err
	}

	// the call and ffmpeg are setup up, we are good to go!

	sm := &statemachine.Statemachine{Name: "Orchestrator"}
	nodeStart := &statemachine.Node{
		Name: "Start",
		OnEnter: func() *statemachine.Node {
			return nil
		},
	}
	sm.AddNode(nodeStart)
	nodeCallEstablished := &statemachine.Node{
		Name: "Call established",
		OnEnter: func() *statemachine.Node {
			log.Printf("Call is established! Playing audio after %d seconds", silenceStartOfCall)
			time.Sleep(time.Duration(silenceStartOfCall) * time.Second)
			close(startPlayFFMPEG)
			return nil
		},
	}
	sm.AddNode(nodeCallEstablished)
	nodeShutdown := &statemachine.Node{
		Name: "Shutdown",
		OnEnter: func() *statemachine.Node {
			close(terminateCall)
			close(terminateFFMPEG)
			return nil
		},
	}
	sm.AddNode(nodeShutdown)

	nodeExit := &statemachine.Node{
		Name: "Exit",
		OnEnter: func() *statemachine.Node {
			return nil
		},
	}
	sm.AddNode(nodeExit)

	sm.Start = nodeStart
	sm.Exit = nodeExit

	sm.AddTransition(&statemachine.Transition{
		From:      nodeStart,
		To:        nodeCallEstablished,
		Condition: callConnected,
	})
	sm.AddTransition(&statemachine.Transition{
		From:      nodeStart,
		To:        nodeShutdown,
		Condition: shutdown,
	})
	sm.AddTransition(&statemachine.Transition{
		From:      nodeStart,
		To:        nodeShutdown,
		Condition: ffmpegTerminated,
	})
	sm.AddTransition(&statemachine.Transition{
		From:      nodeStart,
		To:        nodeShutdown,
		Condition: callTerminated,
	})
	sm.AddTransition(&statemachine.Transition{
		From:      nodeShutdown,
		To:        nodeExit,
		Condition: statemachine.ALWAYS,
	})
	sm.AddTransition(&statemachine.Transition{
		From:      nodeCallEstablished,
		To:        nodeShutdown,
		Condition: ffmpegTerminated,
	})
	sm.AddTransition(&statemachine.Transition{
		From:      nodeCallEstablished,
		To:        nodeShutdown,
		Condition: callTerminated,
	})
	sm.AddTransition(&statemachine.Transition{
		From:      nodeCallEstablished,
		To:        nodeShutdown,
		Condition: shutdown,
	})

	sm.Run()
	log.Println("playTrack completed")

	return nil
}
