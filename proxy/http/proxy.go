package http

import (
	"fmt"
	"github.com/sirupsen/logrus"
	"net/http"
	"time"
)

type HandlerFunc func(proxy *Proxy, res http.ResponseWriter, req *http.Request)
type MethodHandlerMap map[string]HandlerFunc

var DefaultMethodHandlers = make(MethodHandlerMap)

type Proxy struct {
	*http.Server
	MethodHandlers MethodHandlerMap
	RoundTripper http.RoundTripper
}

func NewProxy() *Proxy {
	proxy := &Proxy{
		Server: &http.Server{},
		RoundTripper:   http.DefaultTransport,
		MethodHandlers: DefaultMethodHandlers,
	}
	proxy.Server.Handler = proxy
	return proxy
}

func orPanic(err error, args... interface{}) {
	if err == nil {
		return
	}
	if len(args) == 0 {
		panic(err)
	}

	format := fmt.Sprintf("%s: %%s", args[0])
	newArgs := append(args[1:], err)
	panic(fmt.Sprintf(format, newArgs...))
}

func noPanic(err error) error {
	var recovered error = nil

	if e := recover(); e != nil {
		if er, ok := e.(error); ok {
			recovered = er
		}
	}

	if err != nil && recovered != nil {
		return fmt.Errorf("%s, %s", err, recovered)
	}
	if err != nil {
		return err
	}
	return recovered
}

func noPanicHandle(err error, handle func(error)) {
	err = noPanic(err)
	if err != nil {
		handle(err)
	}
}

func (p *Proxy) ServeHTTP(writer http.ResponseWriter, request *http.Request) {
	logrus.Debugf("ServeHTTP: %v\n", request)
	defer logTime(logrus.Debugf, "end ServeHTTP: %v\n", time.Now())

	handlerFunc := p.MethodHandlers[request.Method]
	if handlerFunc == nil {
		panic("unsupported request method: "+request.Method)
	}

	handlerFunc(p, writer, request)
}

