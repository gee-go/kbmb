package main

import (
	"log"
	"os"
	"os/exec"

	. "github.com/visionmedia/go-gracefully"
)

type Options struct {
	Count int `short:"n" default:"1" description:"number of nsqd nodes"`
}

func check(err error) {
	if err != nil {
		panic(err)
	}
}

func run(quit chan bool, cmd string, args ...string) {
	log.Printf("exec %s %v", cmd, args)
	proc := exec.Command(cmd, args...)
	proc.Stderr = os.Stderr
	proc.Stdout = os.Stdout

	err := proc.Start()
	check(err)

	<-quit

	log.Printf("kill %s", cmd)
	err = proc.Process.Kill()
	check(err)
}

func main() {
	quit := make(chan bool)

	check(os.MkdirAll("/tmp/nsqd", 0755))

	go run(quit, "nsqlookupd")
	go run(quit, "nsqadmin", "--lookupd-http-address", "127.0.0.1:4161")

	go run(quit, "nsqd", "--lookupd-tcp-address", "127.0.0.1:4160", "--data-path", "/tmp/nsqd")
	Shutdown()
	close(quit)
}
