package wxws

type Packet struct {
	Event     Event
	Interests []string
}

type Event struct {
	Name string
	Data string
}
