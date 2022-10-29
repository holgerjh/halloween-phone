package statemachine

// Simple finite state machine implementation

import (
	"fmt"
	"log"
	"reflect"
)

type Transition struct {
	From      *Node
	To        *Node
	Condition <-chan int
}

type Node struct {
	OnEnter func() *Node
	Name    string
}

type Statemachine struct {
	Transitions []*Transition
	Nodes       []*Node
	Start       *Node
	Exit        *Node
	LookupMap   FastLookup
	Name        string
}

var ALWAYS chan int

func init() {
	ALWAYS = make(chan int)
	close(ALWAYS)
}

type FastLookup map[string][]*Transition

func (s *Statemachine) cleanup() {
	s.LookupMap = nil
	for _, v := range s.Transitions {
		v.From = nil
		v.To = nil
	}
	s.Transitions = nil
	s.Nodes = nil
}

func (s *Statemachine) Run() error {
	s.GenerateLookupMap()
	currentNode := s.Start
	if s.Exit == nil {
		return fmt.Errorf("no end node set")
	}
	if currentNode == nil {
		return fmt.Errorf("no start node set")
	}
	for {
		if currentNode.Name == s.Exit.Name {
			log.Printf("[Statemachine %s] Reached exit node, returning", s.Name)
			s.cleanup()
			return nil
		}
		log.Printf("[Statemachine %s] Entering node %s", s.Name, currentNode.Name)
		warp := currentNode.OnEnter()
		if warp != nil {
			log.Printf("[Statemachine %s] Current node requested warp to %s", s.Name, warp.Name)
			currentNode = warp
			continue
		}

		transitions := s.LookupMap[currentNode.Name]
		log.Printf("[Statemachine %s] Found %d outgoing transitions", s.Name, len(transitions))
		cases := make([]reflect.SelectCase, len(transitions))
		for i, v := range transitions {
			cases[i] = reflect.SelectCase{Dir: reflect.SelectRecv, Chan: reflect.ValueOf(v.Condition)}
		}
		chosen, _, recvOK := reflect.Select(cases)
		if recvOK {
			return fmt.Errorf("failed running select statement with cases %+v (channel not closed, make sure to communicate via closing channels!)", cases)
		}
		nextNode := transitions[chosen].To
		log.Printf("[Statemachine %s] Condition %d fulfilled, moving to node %s", s.Name, chosen, nextNode.Name)
		currentNode = nextNode
	}
}

func (s *Statemachine) AddNode(n *Node) {
	s.Nodes = append(s.Nodes, n)
}

func (s *Statemachine) AddTransition(t *Transition) error {
	s.Transitions = append(s.Transitions, t)
	return nil
}

func (s *Statemachine) GenerateLookupMap() {
	s.LookupMap = make(FastLookup)
	for _, node := range s.Nodes {
		if _, ok := s.LookupMap[node.Name]; !ok {
			s.LookupMap[node.Name] = make([]*Transition, 0)
		}
		for _, v := range s.Transitions {
			if v.From.Name == node.Name {
				s.LookupMap[node.Name] = append(s.LookupMap[node.Name], v)
			}
		}
	}
}
