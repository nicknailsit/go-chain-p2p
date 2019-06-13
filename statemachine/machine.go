package statemachine


type SMInterface interface{

	Initial(name string) *StateMachine
	State(name string) *State
	Event(name string) *Event
	Trigger(name string, value SMGSInterface, st *storage, notes ...string) error
}

type SMGSInterface interface {
	SetState(name string)
	GetStateByName(name string) (code byte)
	GetStateByByte(code byte) (name string)
	GetState() string
}

type StateEventInterface interface {
	To(name string) *EventTransition
	From(states ...string) *EventTransition
	Before(func(value interface{}, st *storage) error) *EventTransition
	After(func(value interface{}, st *storage) error) *EventTransition
	Enter(func(value interface{}, st *storage) error) *State
	Exit(func(value interface{}, st *storage) error) *State
	enters(func(value interface{}, st *storage) error) *State
	exits(func(value interface{}, st *storage) error) *State
}

type Transition struct {
	State string
	StateCode byte
	PrevStates [][]byte
	lastTransition int64
}

type StateMachine struct {
	initialState string
	states map[string]*State
	events map[string]*Event
}


type State struct {
	Name string
	Code byte
	enters []func(value interface{}, st *storage) error
	exits []func(value interface{}, st *storage) error
}

type Event struct {
	Name string
	transitions []*EventTransition
}

type EventTransition struct {
	to string
	froms []string
	befores []func(value interface{}, st *storage) error
	afters []func(value interface{}, st *storage) error
}

