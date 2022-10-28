package player

import (
	"bufio"
	"io"
	"log"
	"strings"
)

type CallState struct {
	Ready       <-chan int
	Established <-chan int
	Ended       <-chan int
}

const MARKER_READY = "200 OK ()"
const MARKER_ESTABLISHED = "Call established"
const MARKER_ENDED = "session closed"

func monitorCall(reader io.Reader) *CallState {
	//TODO: check if we need better error handling logic that terminates
	// this go routine
	scanner := bufio.NewScanner(reader)

	callReady := make(chan int)
	callEstablished := make(chan int)
	callEnded := make(chan int)

	go func() {
		log.Printf("Started monitoring of baresip output")
		for scanner.Scan() {
			line := scanner.Text()
			log.Printf("[MONITOR] RAW: %s", line)
			if strings.Contains(line, MARKER_READY) {
				log.Printf("Monitor: Found MARKER_READY")
				close(callReady)
				continue
			}
			if strings.Contains(line, MARKER_ESTABLISHED) {
				log.Printf("Monitor: Found MARKER_ESTALISHED")
				close(callEstablished)
				continue
			}
			if strings.Contains(line, MARKER_ENDED) {
				log.Printf("Monitor: Found MARKER_ENDED")
				close(callEnded)
				continue
			}
		}
	}()
	return &CallState{Ready: callReady, Established: callEstablished, Ended: callEnded}
}
