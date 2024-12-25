package main

import (
	"context"
	"fmt"
	"time"

	"dubbo.apache.org/dubbo-go/v3"
	"dubbo.apache.org/dubbo-go/v3/client"
	"dubbo.apache.org/dubbo-go/v3/registry"
	greet "github.com/zhangcunli/triple-stream-test/proto"
)

func testReuseCli() {
	cli, _ := getClient()
	for i := 0; i < 200; i++ {
		TestBiDiStream(cli)
		time.Sleep(20 * time.Millisecond)
	}

	time.Sleep(3 * time.Second)
	for i := 0; i < 200; i++ {
		TestBiDiStream(cli)
		time.Sleep(20 * time.Millisecond)
	}

	time.Sleep(3 * time.Second)
	for i := 0; i < 200; i++ {
		TestBiDiStream(cli)
		time.Sleep(20 * time.Millisecond)
	}

	time.Sleep(3 * time.Second)
	for i := 0; i < 200; i++ {
		TestBiDiStream(cli)
		time.Sleep(20 * time.Millisecond)
	}

	time.Sleep(3 * time.Second)
	for i := 0; i < 200; i++ {
		TestBiDiStream(cli)
		time.Sleep(20 * time.Millisecond)
	}
	fmt.Printf("=====================================\n")

	//每一次调用 TestBiDiStream 后，都关闭了 stream, 并 sleep
	//但是后续的 TestBiDiStream 没有复用连接池中的 TCP 连接
	//共请求了 1000 次，正好创建了 1000 个 TCP 连接
	//#netstat -an|grep 20000|wc -l
	// 2001
}

func getClient() (*client.Client, error) {
	ins, err := dubbo.NewInstance(
		dubbo.WithName("dubbo_registry_zookeeper_client"),
		dubbo.WithRegistry(
			registry.WithZookeeper(),
			registry.WithAddress("127.0.0.1:2181"),
		),
	)
	if err != nil {
		panic(err)
	}

	cli, err := ins.NewClient()
	if err != nil {
		panic(err)
	}

	return cli, nil
}

func TestBiDiStream(cli *client.Client) error {
	svc, err := greet.NewGreetService(cli)
	if err != nil {
		panic(err)
	}

	fmt.Printf("start to test TRIPLE bidi stream\n")
	stream, err := svc.GreetStream(context.Background())
	if err != nil {
		return err
	}
	if sendErr := stream.Send(&greet.GreetStreamRequest{Name: "Hello! This is triple stream client!"}); sendErr != nil {
		return err
	}

	resp, err := stream.Recv()
	if err != nil {
		return err
	}
	fmt.Printf("TRIPLE bidi stream resp: %s\n", resp.Greeting)
	if err := stream.CloseRequest(); err != nil {
		return err
	}
	if err := stream.CloseResponse(); err != nil {
		return err
	}
	fmt.Printf(">>>>>>>>TestBiDiStream end, close stream...\n")
	return nil
}
