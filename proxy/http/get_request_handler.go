package http

import (
	"bytes"
	"github.com/sirupsen/logrus"
	"io"
	_ "io"
	"io/ioutil"
	"net/http"
	"time"
)

func init() {
	DefaultMethodHandlers[http.MethodGet] = HandleGet
}
type logFunc func(string, ...interface{})

func logTime(logFunc logFunc, format string, start time.Time) {
	logFunc(format, time.Now().Sub(start))
}

func HandleGet(_ *Proxy, w http.ResponseWriter, req *http.Request) {
	logrus.Debugf("GET %v\n", req)
	defer logTime(logrus.Debugf, "GET time: %v\n", time.Now())

	var bodyReader io.Reader = nil
	var bodyBytes []byte = nil

	if req.Body != nil {
		var err error
		bodyBytes, err = ioutil.ReadAll(req.Body)
		orPanic(err, "read body")

		bodyReader = bytes.NewBuffer(bodyBytes)
	}
	newReq, err := http.NewRequest(req.Method, req.URL.String(), bodyReader)
	orPanic(err, "new http request(%s, $s, %v)", req.Method, req.URL.String(), bodyReader)

	res, err := http.DefaultTransport.RoundTrip(newReq)
	if err != nil {
		http.Error(w, err.Error(), http.StatusServiceUnavailable)
		return
	}
	defer res.Body.Close()
	copyHeader(w.Header(), res.Header)
	w.WriteHeader(res.StatusCode)
	io.Copy(w, res.Body)
}

	func copyHeader(src , dst http.Header) {
    for k, vv := range src {
        for _, v := range vv {
            dst.Add(k, v)
        }
    }
}
