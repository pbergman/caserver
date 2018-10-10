package router

import (
	"fmt"
	"net/http"
	"strconv"
	"sync"

	"github.com/pbergman/logger"
)

type Router struct {
	controllers []ControllerInterface
	pre         []PreControllerInterface
	logger      *logger.Logger
	lock        sync.RWMutex
}

type writeWrapper struct {
	http.ResponseWriter
	statusCode int
}

func newWriteWrapper(inner http.ResponseWriter) http.ResponseWriter {
	return &writeWrapper{inner, http.StatusOK}
}

func (r *writeWrapper) WriteHeader(code int) {
	r.statusCode = code
	r.ResponseWriter.WriteHeader(code)
}

func (r writeWrapper) statusLine() string {
	return strconv.Itoa(r.statusCode) + " " + http.StatusText(r.statusCode)
}

func (r *Router) ServeHTTP(response http.ResponseWriter, request *http.Request) {
	r.logger.Debug(r.requestLine(request))
	response = newWriteWrapper(response)
	//defer runtime.GC()

	defer func() {
		r.logger.Debug(fmt.Sprintf("[%p] %s", request, response.(*writeWrapper).statusLine()))
	}()

	if request.RequestURI == "*" {
		if request.ProtoAtLeast(1, 1) {
			response.Header().Set("Connection", "close")
		}
		response.WriteHeader(http.StatusBadRequest)
		return
	}

	for i, c := 0, len(r.pre); i < c; i++ {
		if r.pre[i].Match(request) {
			r.pre[i].Handle(request, response.Header(), r.logger.Get(r.pre[i].Name()))
		}
	}

	if controller := r.handler(request); controller != nil {
		controller.Handle(response, request, r.logger.Get(controller.Name()))
		return
	}

	http.NotFound(response, request)
}

func (r *Router) handler(request *http.Request) ControllerInterface {
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
	uri := req.URL.Path
	if uri == "" {
		uri = req.URL.RequestURI()
	}
	method := req.Method
	if method == "" {
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
