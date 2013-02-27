package main

import (
	"sync"
)


type State struct {
	Lock sync.Mutex
	State map[*Target]TargetStatus
}

func NewState() *State{
	s := new(State)
	s.State = make(map[*Target]TargetStatus)
	return s
}

