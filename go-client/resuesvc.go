package main

import (
	"context"
	"fmt"
	"time"

	"dubbo.apache.org/dubbo-go/v3"
	"dubbo.apache.org/dubbo-go/v3/registry"
	greet "github.com/zhangcunli/triple-stream-test/proto"
)

func testReuseSvc() {
	svc, _ := getService()
	for i := 0; i < 200; i++ {
		TestBiDiStream2(svc)
		time.Sleep(20 * time.Millisecond)
	}

	time.Sleep(3 * time.Second)
	for i := 0; i < 200; i++ {
		TestBiDiStream2(svc)
		time.Sleep(20 * time.Millisecond)
	}
}

func getService() (greet.GreetService, error) {
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

	svc, err := greet.NewGreetService(cli)
	if err != nil {
		panic(err)
	}
	return svc, nil
}

func TestBiDiStream2(svc greet.GreetService) error {
	fmt.Printf("start to test triple bidi stream 2\n")
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
	fmt.Printf("triple bidi stream2 resp: %s\n", resp.Greeting)
	if err := stream.CloseRequest(); err != nil {
		return err
	}
	if err := stream.CloseResponse(); err != nil {
		return err
	}
	fmt.Printf("========>TestBiDiStream end, close stream...\n")
	return nil
}
