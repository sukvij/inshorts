package interactionservice

import (
	"encoding/json"
	"fmt"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/go-redis/redis"
	"github.com/mmcloughlin/geohash"
	llmservice "github.com/sukvij/inshorts/all-services/llm-service"
	newsservice "github.com/sukvij/inshorts/all-services/news-service"
	"github.com/sukvij/inshorts/inshortfers/logs"
	redisservice "github.com/sukvij/inshorts/inshortfers/redis-service"
	"github.com/sukvij/inshorts/inshortfers/response"
	"gorm.io/gorm"
)

type InteractionController struct {
	DB          *gorm.DB
	Interaction *[]UserInteraction
	Redis       *redis.Client
	Logs        *logs.AgreeGateLoager
}

func _NewController(db *gorm.DB, redis *redis.Client, logs *logs.AgreeGateLoager) *InteractionController {
	return &InteractionController{DB: db, Redis: redis, Logs: logs}
}

func UserInteractionController(appGroup *gin.RouterGroup, db *gorm.DB, redis *redis.Client, logs *logs.AgreeGateLoager) {
	// /v1/interaction
	app := appGroup.Group("/interaction")
	controller := _NewController(db, redis, logs)
	app.POST("", controller.createUserInteraction)
	app.GET("/trending", controller.trendingNewsArticles)
}

func (controller *InteractionController) createUserInteraction(ctx *gin.Context) {
	var interactions []UserInteraction
	err := ctx.ShouldBindJSON(&interactions)
	if err != nil {
		controller.Logs.Error(err)
		response.JSONResponse(ctx, err, nil, nil, "")
		return
	}
	service := _NewService(controller.DB, &interactions, controller.Redis)
	err = service.CreateUserInteraction()
	if err != nil {
		controller.Logs.Error(err)
	}
	response.JSONResponse(ctx, err, nil, controller.Redis, "")
}

func (controller *InteractionController) trendingNewsArticles(ctx *gin.Context) {
	lat := ctx.Query("lat")
	lon := ctx.Query("lon")
	limit := ctx.Query("limit")
	lat1, _ := strconv.ParseFloat(lat, 64)
	lon1, _ := strconv.ParseFloat(lon, 64)
	radiusMeters := 100000
	limit1, _ := strconv.Atoi(limit)
	geohashPrecision := 6
	geoHashKey := geohash.EncodeWithPrecision(lat1, lon1, uint(geohashPrecision))
	cacheKey := fmt.Sprintf("v1:interaction:trending:%s:limit%d:radius%d", geoHashKey, limit1, radiusMeters)
	fmt.Println("key is this bro", cacheKey)

	// check if data is rpesent in cache or not
	val, err := redisservice.GetValueFromRedis(controller.Redis, cacheKey)
	if err == redis.Nil {
		var ans response.FinalResponse
		temp, _ := json.Marshal(val)
		json.Unmarshal(temp, &ans)
		ctx.JSON(200, ans)
		return
	}
	fmt.Println("database me jaa rha h matlab --> redis khali")
	service := _NewService(controller.DB, &[]UserInteraction{}, controller.Redis)
	result, err := service.TrendingNewsArticles(lat1, lon1, limit1, cacheKey, radiusMeters)
	if err != nil {
		controller.Logs.Error(err)
		response.JSONResponse(ctx, err, nil, nil, "")
		return
	}
	ConvertResultInProperFormatAndReturn(ctx, result, controller.Redis, cacheKey)
}

func ConvertResultInProperFormatAndReturn(ctx *gin.Context, result *[]newsservice.NewsArticle, redisClient *redis.Client, cacheKey string) {
	finalResult := newsservice.Convert_NewsArticle_To_NewsArticleResponse(result)
	// generate llm summery
	for i := 0; i < len(*finalResult); i++ {
		summery := llmservice.GenerateSummeryLLM((*finalResult)[i].Title, (*finalResult)[i].Description)
		(*finalResult)[i].LLMSummery = summery
	}
	response.JSONResponse(ctx, nil, finalResult, redisClient, cacheKey)
}
