package main

import (
	"fmt"
	"log"
	"net"

	"github.com/go-redis/redis"
	"github.com/yeralin-munar/codding_interview/cache/handler"
	"github.com/yeralin-munar/codding_interview/cache/utils"
	"github.com/yeralin-munar/codding_interview/generated"
	"google.golang.org/grpc"
)

func main() {
	log.Println("Start cache-service listen on port 7777")
	conf := utils.GetConf("config.yml")
	log.Println(conf)
	rdb := redis.NewClient(&redis.Options{
		Addr:     "redis:6379",
		Password: "", // no password set
		DB:       0,  // use default DB
	})
	redisStatusCmd := rdb.Ping()
	if redisStatusCmd.Err() != nil {
		log.Fatalf("[ERROR] Make connect to Redis: %s", redisStatusCmd.Err().Error())
	}
	log.Println(redisStatusCmd.Val())
	netListener := getNetListener(7777)
	gRPCServer := grpc.NewServer()

	cacheService := handler.NewCacheService(conf, rdb)

	generated.RegisterCacheServer(gRPCServer, cacheService)

	if err := gRPCServer.Serve(netListener); err != nil {
		log.Fatalf("Failed to start cache-service: %s", err)
	}

}

func getNetListener(port uint) net.Listener {
	listener, err := net.Listen("tcp", fmt.Sprintf(":%d", port))
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
		panic(fmt.Sprintf("failed to listen: %v", err))
	}

	return listener
}
