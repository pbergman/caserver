package controller

import (
    "net/http"
    "github.com/pbergman/logger"
)

type PreAccessControl struct {

}

func (p PreAccessControl) Name() string {
    return "pre.hook.access.control.header"
}

func (p PreAccessControl) Match(request *http.Request) bool {
    return request.Header.Get("accept") == "application/json"
}

func (p PreAccessControl) Handle(request *http.Request, header http.Header, logger logger.LoggerInterface) {
    logger.Debug("adding 'Access-Control-Allow-Origin: *' to the response")
    header.Add("access-Control-Allow-Origin", "*")
}