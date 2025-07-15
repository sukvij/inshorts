package interactionservice

import (
	"encoding/json"
	"fmt"
	"strconv"

	"github.com/go-redis/redis"
	"github.com/mmcloughlin/geohash"
	newsservice "github.com/sukvij/inshorts/all-services/news-service"
	redisservice "github.com/sukvij/inshorts/inshortfers/redis-service"
	"gorm.io/gorm"
)

type InteractionService struct {
	DB           *gorm.DB
	Interactions *[]UserInteraction
	Redis        *redis.Client
}

func _NewService(db *gorm.DB, interatcions *[]UserInteraction, redis *redis.Client) *InteractionService {
	return &InteractionService{DB: db, Interactions: interatcions, Redis: redis}
}

func (service *InteractionService) CreateUserInteraction() error {
	repo := _NewRepository(service.DB, service.Interactions)
	return repo.CreateUserInteraction()
}

func (service *InteractionService) TrendingNewsArticles(lat, lon, limit string) (*[]newsservice.NewsArticle, error) {
	lat1, _ := strconv.ParseFloat(lat, 64)
	lon1, _ := strconv.ParseFloat(lon, 64)
	radiusMeters := 100000
	limit1, _ := strconv.Atoi(limit)
	geohashPrecision := 6
	geoHashKey := geohash.EncodeWithPrecision(lat1, lon1, uint(geohashPrecision))
	cacheKey := fmt.Sprintf("trending:%s:limit%d:radius%d", geoHashKey, limit1, radiusMeters)
	fmt.Println("key is this bro", cacheKey)

	// check if data is rpesent in cache or not
	val, err := redisservice.GetValueFromRedis(service.Redis, cacheKey)
	if err == redis.Nil {
		var ans []newsservice.NewsArticle
		temp, _ := json.Marshal(val)
		json.Unmarshal(temp, &ans)
		return &ans, nil
	}
	fmt.Println("database me jaa rha h matlab --> redis khali")
	repo := _NewRepository(service.DB, service.Interactions)
	res, errs := repo.TrendingNewsArticles(lat1, lon1, limit1, radiusMeters)
	if errs == nil {
		redisservice.SetValueToRedis(service.Redis, cacheKey, *res)
	}
	return res, errs
}
