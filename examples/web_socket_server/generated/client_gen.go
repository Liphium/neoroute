// Code generated with neogen-generated v1 schema by neogen. DO NOT EDIT.
package main

type NewPunEvent struct {
	Pun        string       `msg:"pun"`
	AnotherOne *NewPunEvent `msg:"another_one"`
}

type EchoRequest struct {
	Message string `msg:"message"`
}

type EchoResponse struct {
	Message       string `msg:"message"`
	RequestNumber int64  `msg:"requestNumber"`
}

type SubmitPunRequest struct {
	Pun string `msg:"pun"`
}
