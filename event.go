package neoroute

//go:generate msgp -unexported

//msgp:tuple event
type event struct {
	Name string `msg:"name"`
	Data []byte `msg:"data"`
}
