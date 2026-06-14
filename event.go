package neoroute

//go:generate msgp -unexported

type event struct {
	Name string `msg:"name"`
	Data []byte `msg:"data"`
}
