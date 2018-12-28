package executor

import (
	"net"

	"github.com/mostafa-asg/pelican/socket/parser"
	"github.com/mostafa-asg/pelican/store"
)

type PutExecutor struct{}

func (e *PutExecutor) Accept(result parser.ParseResult) bool {
	return result.Command() == "PUT"
}

func (e *PutExecutor) Execute(result parser.ParseResult, store *store.Store, con net.Conn) {
	p := result.(*parser.PutResult)
	store.Put(p.Key, p.Value)
	con.Write([]byte("OK."))
}

func NewPutExecutor() *PutExecutor {
	return &PutExecutor{}
}
