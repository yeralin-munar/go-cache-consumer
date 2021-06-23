package main

import (
	"context"
	"log"
	"sync"

	"github.com/yeralin-munar/codding_interview/generated"
	"google.golang.org/grpc"
	"google.golang.org/protobuf/types/known/emptypb"
)

func main() {
	serverAddress := "localhost:7777"

	conn, e := grpc.Dial(serverAddress, grpc.WithInsecure())

	if e != nil {
		panic(e)
	}
	defer conn.Close()

	client := generated.NewCacheClient(conn)

	var wg sync.WaitGroup
	for i := 0; i < 1000; i++ {
		wg.Add(1)
		go makeRequest(&wg, &client)
	}
	wg.Wait()
}

func makeRequest(wg *sync.WaitGroup, client *generated.CacheClient) {
	defer wg.Done()
	stringStream, err := (*client).GetRandomDataStream(context.Background(), &emptypb.Empty{})
	if err != nil {
		log.Printf("[ERROR] Make request to server: %s", err.Error())
	}
	log.Println(len(stringStream.String()))
}
