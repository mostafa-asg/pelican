package executor

import (
	"net"

	"github.com/mostafa-asg/pelican/socket/parser"
	"github.com/mostafa-asg/pelican/store"
)

type PutEExecutor struct{}

func (e *PutEExecutor) Accept(result parser.ParseResult) bool {
	return result.Command() == "PUTE"
}

func (e *PutEExecutor) Execute(result parser.ParseResult, store *store.Store, con net.Conn) {
	p := result.(*parser.PutEResult)
	store.PutWithExpire(p.Key, p.Value, p.Expiration, p.Strategy)
	con.Write([]byte("OK."))
}

func NewPutEExecutor() *PutEExecutor {
	return &PutEExecutor{}
}
