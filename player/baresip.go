package player

import (
	"github.com/holgerjh/halloween-phone/statemachine"
	"io"
	"log"
	"os"
	"os/exec"
	"time"
)

const bareSipBinary = "baresip"

func Call(maxRingTime int, cancelCall <-chan int) (<-chan int, <-chan int, error) {
	callConnected := make(chan int)
	callTerminated := make(chan int)

	log.Printf("Establishing call")
	bareSipPath, err := exec.LookPath(bareSipBinary)
	if err != nil {
		return nil, nil, err
	}
	go doCall(bareSipPath, maxRingTime, callConnected, callTerminated, cancelCall)

	log.Printf("Spawned calling goroutine, returning channels")
	return callConnected, callTerminated, nil

}

func doCall(bareSipPath string, maxRingTime int, callConnected, callTerminated chan int, cancelCall <-chan int) {
	log.Printf("[Caller] Creating pipes")
	stdinRead, stdinWrite := io.Pipe()
	stdoutRead, stdoutWrite := io.Pipe()

	log.Printf("[Caller] Creating baresip command and wiring pipes")
	cmd := exec.Command(bareSipPath)
	cmd.Stdin = stdinRead
	cmd.Stdout = stdoutWrite
	cmd.Stderr = os.Stderr

	log.Printf("[Caller] Spawning goroutine running %s", bareSipPath)
	baresipTerminated := make(chan int)
	go func() {
		err := cmd.Run()
		if err != nil {
			log.Printf("Baresip terminated with error: %e", err)
		}
		log.Printf("baresip terminated, closing channel")
		close(baresipTerminated)
	}()

	callstate := monitorCall(stdoutRead)

	log.Printf("Building state machine")
	sm := &statemachine.Statemachine{Name: "baresip"}

	nodeExit := &statemachine.Node{
		Name: "Exit node",
		OnEnter: func() *statemachine.Node {
			return nil
		},
	}
	sm.AddNode(nodeExit)
	sm.Exit = nodeExit

	// shutdown process
	nodeShutdownProcess1 := &statemachine.Node{
		Name: "Shutdown baresip",
		OnEnter: func() *statemachine.Node {
			log.Printf("Terminating baresip")
			terminateBaresip(stdinWrite, cmd)
			return nil
		},
	}
	sm.AddNode(nodeShutdownProcess1)
	nodeShutdownProcess2 := &statemachine.Node{
		Name: "Inform about terminated call",
		OnEnter: func() *statemachine.Node {
			log.Printf("Informing about terminated call")
			close(callTerminated)
			return nil
		},
	}
	sm.AddNode(nodeShutdownProcess2)
	sm.AddTransition(
		&statemachine.Transition{
			From:      nodeShutdownProcess1,
			Condition: baresipTerminated,
			To:        nodeShutdownProcess2})
	sm.AddTransition(
		&statemachine.Transition{
			From:      nodeShutdownProcess2,
			Condition: statemachine.ALWAYS,
			To:        nodeExit})

	// wait ready -> is ready -> established -> done
	nodeWaitReadyForCall := &statemachine.Node{
		Name:    "Waiting for established Call",
		OnEnter: func() *statemachine.Node { return nil },
	}
	sm.AddNode(nodeWaitReadyForCall)
	nodeIsReadyForCall := &statemachine.Node{
		Name: "IsReadyFor Call",
		OnEnter: func() *statemachine.Node {
			log.Printf("Dialing")
			if _, err := stdinWrite.Write([]byte("/dialcontact\n")); err != nil {
				log.Printf("ERROR: Unable to communicate with baresip: %e", err)
				return nodeShutdownProcess1
			}
			return nil
		},
	}
	sm.AddNode(nodeIsReadyForCall)
	nodeCallEstablished := &statemachine.Node{
		Name: "Established Call",
		OnEnter: func() *statemachine.Node {
			log.Printf("Someone picked up!")
			close(callConnected)
			return nil
		},
	}
	sm.AddNode(nodeCallEstablished)
	sm.AddTransition(
		&statemachine.Transition{
			From:      nodeWaitReadyForCall,
			Condition: cancelCall,
			To:        nodeShutdownProcess1,
		})
	sm.AddTransition(
		&statemachine.Transition{
			From:      nodeWaitReadyForCall,
			Condition: callstate.Ready,
			To:        nodeIsReadyForCall,
		})

	sm.AddTransition(
		&statemachine.Transition{
			From:      nodeIsReadyForCall,
			Condition: cancelCall,
			To:        nodeShutdownProcess1,
		})
	sm.AddTransition(
		&statemachine.Transition{
			From:      nodeIsReadyForCall,
			Condition: callstate.Established,
			To:        nodeCallEstablished,
		})

	sm.AddTransition(
		&statemachine.Transition{
			From:      nodeCallEstablished,
			Condition: cancelCall,
			To:        nodeShutdownProcess1,
		})
	sm.AddTransition(
		&statemachine.Transition{
			From:      nodeCallEstablished,
			Condition: callstate.Ended,
			To:        nodeShutdownProcess1,
		})

	sm.Start = nodeWaitReadyForCall

	log.Printf("Running state machine")
	if err := sm.Run(); err != nil {
		panic(err)
	}

	/*
	   select {
	   // happy path: ready -> established -> ended -> return
	   case <-cancelCall:

	   	log.Printf("Received request to terminate call, terminating baresip")
	   	terminateBaresip(stdinWrite, cmd)

	   case <-callstate.Ready:

	   		log.Printf("Dialing")
	   		if _, err := stdinWrite.Write([]byte("/dialcontact\n")); err != nil {
	   			log.Printf("ERROR: Unable to communicate with baresip: %e", err)
	   			terminateBaresip(stdinWrite, cmd)
	   		}
	   	}

	   case <-callstate.Established:

	   	log.Printf("Someone picked up!")
	   	close(callConnected)

	   case <-callstate.Ended:

	   	log.Printf("call completed, terminating baresip")
	   	terminateBaresip(stdinWrite, cmd)

	   // baresip shutdown: signal end of call and leave function
	   case <-baresipTerminated:

	   	log.Printf("baresip terminated, terminating call")
	   	close(callTerminated)
	   	return

	   // special case: call timeout
	   case <-time.After(time.Duration(maxRingTime) * time.Second):

	   	log.Printf("Noone picked up, terminating baresip")
	   	terminateBaresip(stdinWrite, cmd)

	   // special case: external request to stop call
	   case <-cancelCall:

	   		log.Printf("Received request to terminate call, terminating baresip")
	   		terminateBaresip(stdinWrite, cmd)
	   	}
	*/
}

func terminateBaresip(stdinWrite *io.PipeWriter, cmd *exec.Cmd) {
	stdinWrite.Close()
	tryKillProcess(cmd)
	/*
		if _, err := stdinWrite.Write([]byte("/quit\n")); err != nil {
			log.Printf("ERROR: Unable to communicate with baresip: %e", err)
			tryKillProcess(cmd)
		} else {
			return //TODO: make sure to kill process after timeout if it ignores /quit command
		}*/
}

const KILL_WAIT_FOR_PROCESS_PTR = 5

func tryKillProcess(cmd *exec.Cmd) {
	log.Printf("Terminating process %s", cmd.Path)
	if cmd.Process == nil {
		log.Printf("WARN: Unable to terminate process: nil pointer to process")
		log.Printf("Waiting %d seconds and try to kill again", KILL_WAIT_FOR_PROCESS_PTR)
		time.Sleep(time.Duration(KILL_WAIT_FOR_PROCESS_PTR) * time.Second)
		if cmd.Process == nil {
			log.Printf("ERROR: unable to terminate process %s due to nil pointer", cmd.Path)
			return
		}
		log.Printf("Process became available, terminating")
	}
	if err := cmd.Process.Kill(); err != nil {
		log.Printf("ERROR: Failed terminating process: %e", err)
	}
}

/*

func waitForEstablishedCall(reader io.Reader, established chan<- int) {
	scanner := bufio.NewScanner(reader)
	for scanner.Scan() {
		if strings.Contains(scanner.Text(), MARKER_ESTABLISHED) {
			close(established)
			return
		}
	}

}*/
