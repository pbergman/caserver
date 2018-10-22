package router

import (
	"fmt"
	"net/http"

	"github.com/pbergman/logger"
)

type Router struct {
	controllers []ControllerInterface
	pre         []PreControllerInterface
	logger      *logger.Logger
}

func (r *Router) ServeHTTP(response http.ResponseWriter, request *http.Request) {

	r.logger.Debug(r.requestLine(request))
	response = newWrappedResponse(response)

	defer func() {
		r.logger.Debug(fmt.Sprintf("[%p] %s", request, response.(*wrappedResponse).statusLine()))
	}()

	if request.RequestURI == "*" {
		if request.ProtoAtLeast(1, 1) {
			response.Header().Set("Connection", "close")
		}
		response.WriteHeader(http.StatusBadRequest)
		return
	}

	wrapped := &Request{Request: request}

	for i, c := 0, len(r.pre); i < c; i++ {
		if r.pre[i].Match(wrapped) {
			r.pre[i].Handle(wrapped, response.Header(), r.logger.Get(r.pre[i].Name()))
		}
	}

	if handler := r.getHandler(wrapped); handler != nil {
		handler.Handle(wrapped, response, r.logger.Get(handler.Name()))
		return
	}

	http.NotFound(response, request)
}

func (r *Router) getHandler(request *Request) ControllerInterface {
	for i, c := 0, len(r.controllers); i < c; i++ {
		if r.controllers[i].Match(request) {
			return r.controllers[i]
		}
	}
	return nil
}

func (r *Router) AddPreHook(hook PreControllerInterface) {
	r.pre = append(r.pre, hook)
}

func (r *Router) requestLine(req *http.Request) string {
	var uri, method string
	if uri = req.URL.Path; uri == "" {
		uri = req.URL.RequestURI()
	}
	if method = req.Method; method == "" {
		method = "GET"
	}
	return fmt.Sprintf("[%p] %s %s HTTP/%d.%d", req, method, uri, req.ProtoMajor, req.ProtoMinor)
}

func NewRouter(logger *logger.Logger, controller ...ControllerInterface) *Router {
	return &Router{
		controllers: controller,
		pre:         make([]PreControllerInterface, 0),
		logger:      logger,
	}
}
