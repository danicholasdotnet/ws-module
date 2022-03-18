package wxws

type Handlers = map[string]Handler

type Handler func(*Packet, *Connection) (Message, error)

func NewHandlers() Handlers {
	return Handlers{
		"pageChange": populationHandler,
	}
}

var HandlersArray = NewHandlers()
