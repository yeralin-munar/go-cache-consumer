package handler

import (
	"bytes"
	"context"
	"log"
	"math/rand"
	"time"

	"github.com/go-redis/redis"
	"github.com/yeralin-munar/codding_interview/cache/utils"
	"github.com/yeralin-munar/codding_interview/generated"
	"google.golang.org/protobuf/types/known/emptypb"
)

type CacheService struct {
	config *utils.Conf
	rdb    *redis.Client
}

func NewCacheService(conf *utils.Conf, rdb *redis.Client) *CacheService {
	return &CacheService{
		config: conf,
		rdb:    rdb}
}

func (c *CacheService) GetRandomDataStream(ctx context.Context, in *emptypb.Empty) (*generated.StringStreamResponse, error) {
	response := &generated.StringStreamResponse{}

	streamChan := make(chan string)

	for i := 0; i < c.config.NumberOfRequests; i++ {
		now := time.Now().UnixNano()
		ctxWithValue := context.WithValue(ctx, "id", now)
		go c.getData(ctxWithValue, streamChan)
	}
	var streamBuf bytes.Buffer
	for i := 0; i < c.config.NumberOfRequests; i++ {
		streamBuf.WriteString(<-streamChan)
	}
	response.Data = streamBuf.String()
	return response, nil
}

func (c *CacheService) getData(ctx context.Context, streamChan chan string) {
	url := c.getRandomURL()
	// Add id to check goroutine work
	ID := ctx.Value("id").(int64)
	log.Printf("Goroutine ID: %d", ID)
	var isWaiting bool
	for true {
		urlValue, err := c.rdb.Get(url).Result()
		if err != nil && err != redis.Nil {
			log.Printf("[ERROR][%d] Error to get data from Redis by key %s: %s", ID, url, err.Error())
		}
		if urlValue != "" {
			if urlValue == "wait" {
				if !isWaiting {
					log.Printf("[%d] Waiting for data from Redis by key %s", ID, url)
					isWaiting = true
				}
				time.Sleep(10 * time.Millisecond)
				continue
			}
			log.Printf("[%d] Get data from Redis by key: %s", ID, url)
			streamChan <- urlValue
			return
		} else {
			log.Printf("[%d] Set wait for key: %s", ID, url)
			redisStatusCmd := c.rdb.Set(url, "wait", 0)
			if redisStatusCmd.Err() != nil {
				log.Printf("[ERROR][%d] Set key %s to Redus: %s", ID, url, err)
			}
			break
		}

	}
	urlValue, err := utils.MakeRequest(url)
	if err != nil {
		log.Printf("[ERROR] %s", err)
	}
	log.Printf("[%d] Get data from URL: %s", ID, url)
	streamChan <- urlValue

	ttl := c.getRandomTTL()
	log.Printf("[%d] Set data from URL to Redis: %s", ID, url)
	redisStatusCmd := c.rdb.Set(url, urlValue, ttl*time.Millisecond)
	if redisStatusCmd.Err() != nil {
		log.Printf("[ERROR][%d] Set key %s to Redus: %s", ID, url, err)
	}
}

func (c *CacheService) getRandomURL() string {
	rand.Seed(int64(rand.Int()))
	randURLIndex := rand.Intn(len(c.config.URLList))
	return c.config.URLList[randURLIndex]
}

func (c *CacheService) getRandomTTL() time.Duration {
	rand.Seed(int64(rand.Int()))
	ttl := rand.Intn(c.config.MaxTimeout-c.config.MinTimeout) + c.config.MinTimeout
	log.Printf("Random TTL: %d", ttl)
	return time.Duration(ttl)
}
