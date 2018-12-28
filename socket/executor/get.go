package executor

import (
	"net"

	"github.com/mostafa-asg/pelican/socket/parser"
	"github.com/mostafa-asg/pelican/store"
	bytesUtil "github.com/mostafa-asg/pelican/util/bytes"
)

type GetExecutor struct{}

func (e *GetExecutor) Accept(result parser.ParseResult) bool {
	return result.Command() == "GET"
}

func (e *GetExecutor) Execute(result parser.ParseResult, store *store.Store, con net.Conn) {
	g := result.(*parser.GetResult)
	value, found := store.Get(g.Key)
	if found {
		con.Write(bytesUtil.ToBytes(uint32(len(value.([]byte)))))
		con.Write(value.([]byte))
	} else {
		con.Write(make([]byte, 4))
	}
}

func NewGetExecutor() *GetExecutor {
	return &GetExecutor{}
}
