// Code generated with neogen-generated v1 schema by neogen. DO NOT EDIT.
package main

type EchoRequest struct {
	Message string `msg:"message"`
}

type EchoResponse struct {
	RequestNumber int32  `msg:"requestNumber"`
	Message       string `msg:"message"`
}

type SubmitPunRequest struct {
	Pun string `msg:"pun"`
}

type NewPunEvent struct {
	AnotherOne *NewPunEvent `msg:"another_one"`
	Pun        string       `msg:"pun"`
}
