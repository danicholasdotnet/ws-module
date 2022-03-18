package wxws

type Handler func(*Packet, *Connection) (Message, error)

type Handlers = map[string]Handler

func NewHandlers() Handlers {
	return Handlers{
		"pageChange": populationHandler,
	}
}

var Population = make(map[int]string, 100)

func populationHandler(p *Packet, c *Connection) (Message, error) {
	var m Message
	var e error
	Population[c.ID] = p.Event.Data
	m.Data = Population
	m.Channel = "population"
	return m, e
}
