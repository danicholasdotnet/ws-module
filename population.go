package wxws

var Population = make(map[int]string, 100)

func populationHandler(p *Packet, c *Connection) (Message, error) {
	var m Message
	var e error
	Population[c.ID] = p.Event.Data
	m.Data = Population
	m.Channel = "population"
	return m, e
}
