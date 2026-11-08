package client

//go:generate msgp -unexported

// Incomming messages

const (
	MessageTypeResponse = 0
	MessageTypeEvent    = 1
)

//msgp:tuple message
type message struct {
	Type int    `msg:"type"` // Response or event
	Data []byte `msg:"data"`
}

//msgp:tuple event
type event struct {
	Name string `msg:"name"`
	Data []byte `msg:"data"`
}

//msgp:tuple response
type response struct {
	Id      int    `msg:"id"`
	HasData bool   `msg:"has_data"`
	IsError bool   `msg:"error"`
	Data    []byte `msg:"data"`
}

// Outgoing messages

//msgp:tuple request
type request struct {
	Id    int    `msg:"id"`
	Route string `msg:"route"`
	Data  []byte `msg:"data"`
}
