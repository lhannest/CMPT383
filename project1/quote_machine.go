package main

type state_interface interface {
	evalute(i int, s string) bool
	getNextStates() []state
	isFinal() bool
}

type state struct {
	name string
	next_states [] state
	is_final bool
}

func (self *state) evalute(i int, s string) bool {
	return false
}

func (self *state) getNextStates() []state {
	return self.next_states
}

func (self *state) isFinal() bool {
	return self.is_final
}

type quote_state struct {
	next_states []state
	is_final bool
}

func (self *quote_state) evalute(i int, s string) bool {
	if s[i] == '"' {
		return true
	} else {
		return false
	}
}