package router

import (
	"net/http"

	"github.com/pbergman/caserver/util"
)

type Request struct {
	*http.Request
	accept *AcceptResponses
}

func (r *Request) GetAcceptResponseType() *AcceptResponses {
	if nil == r.accept {
		r.accept = NewAcceptResponses(util.DefaultString(r.Header.Get("accept"), "text/plain"))
	}
	return r.accept
}
