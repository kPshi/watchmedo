package main

import (
	"./net"
	httpProxy "./proxy/http"
	"flag"
	"fmt"
	log "github.com/sirupsen/logrus"
	"os"
	"os/exec"
	"syscall"
)


func main() {
	log.SetLevel(log.DebugLevel)
	fh, err := os.OpenFile("watchmedo.log", os.O_WRONLY | os.O_APPEND | os.O_CREATE, 0644)
	orPanic(err, "write to watchmedo.log")
	log.SetOutput(fh)
	httpProxyAddress := flag.String("http-proxy-address", "127.0.0.1:0", "HTTP proxy listen address")
	listener, err := net.Listen("tcp",*httpProxyAddress)
	orPanic(err, "listen")

	log.Printf("listening on %s", listener.Addr().String())

	prx := httpProxy.NewProxy()
	go func() {
		log.Fatalf("Failed to serve: %s\n", prx.Serve(listener))
	}()

	cmd := []string{os.Getenv("SHELL")}
	if len(os.Args)>1 {
		cmd = os.Args[1:]
	}
	log.Infof("starting %v", cmd)
	proxyAddress := fmt.Sprintf("http://%s", listener.Addr().String())
	proc, err := startCommand(cmd, []string{
		"http_proxy="+proxyAddress,
		"https_proxy="+proxyAddress,
		"HTTP_PROXY="+proxyAddress,
		"HTTPS_PROXY="+proxyAddress,
	})
	orPanic(err, "start command")

	log.Infof("Waiting for process to exit.")
	_, err = proc.Wait()
	if err != nil {
		if exitErr, ok := err.(*exec.ExitError); ok {
			if status, ok := exitErr.Sys().(syscall.WaitStatus); ok {
				os.Exit(status.ExitStatus())
			}
		}
	}
	orPanic(err, "wait")
	os.Exit(0)
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

func startCommand(cmd []string, env []string) (*os.Process, error) {
	cwd, err := os.Getwd()
	if err != nil {
		return nil, fmt.Errorf("get cwd: %s", err)
	}

	procAttr := &os.ProcAttr{
		Files: []*os.File{
			os.Stdin,
			os.Stdout,
			os.Stderr,
		},
		Dir: cwd,
		Env: env,
	}

	return os.StartProcess(cmd[0], cmd, procAttr)
}
