package executor

import (
	"net"

	"github.com/mostafa-asg/pelican/socket/parser"
	"github.com/mostafa-asg/pelican/store"
)

type DelExecutor struct{}

func (e *DelExecutor) Accept(result parser.ParseResult) bool {
	return result.Command() == "DEL"
}

func (e *DelExecutor) Execute(result parser.ParseResult, store *store.Store, con net.Conn) {
	d := result.(*parser.DelResult)
	store.Del(d.Key)
	con.Write([]byte("OK."))
}

func NewDelExecutor() *DelExecutor {
	return &DelExecutor{}
}
