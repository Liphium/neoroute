package neoroute

//go:generate msgp -unexported

type request struct {
	Id    int    `msg:"id"`
	Route string `msg:"route"`
	Data  []byte `msg:"data"`
}
