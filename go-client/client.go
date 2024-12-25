package main

import (
	"fmt"
	"os"
	"os/signal"

	_ "dubbo.apache.org/dubbo-go/v3/imports"
)

func main() {
	//testReuseCli()
	testReuseSvc()

	//testReuseCli：复用 client(ins.NewClient())，每次都会建新的 TCP 连接；
	//              用户自己再维护一个 client pool?
	//testReuseSvc：复用 service(greet.NewGreetService(cli))，只会创建一个 TCP 连接

	errc := make(chan error, 1)
	sigs := make(chan os.Signal, 1)
	signal.Notify(sigs, os.Interrupt)
	select {
	case err := <-errc:
		fmt.Printf("failed to serve: %v", err)
	case sig := <-sigs:
		fmt.Printf("terminating: %v", sig)
	}
}
