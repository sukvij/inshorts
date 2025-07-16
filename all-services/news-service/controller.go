package newsservice

import (
	"encoding/json"
	"fmt"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/go-redis/redis"
	llmservice "github.com/sukvij/inshorts/all-services/llm-service"
	redisservice "github.com/sukvij/inshorts/inshortfers/redis-service"
	"github.com/sukvij/inshorts/inshortfers/response"
	"gorm.io/gorm"
)

type NewsController struct {
	DB          *gorm.DB
	NewsArticle []*NewsArticle
	Redis       *redis.Client
}

func _NewController(db *gorm.DB, redis *redis.Client) *NewsController {
	return &NewsController{DB: db, Redis: redis}
}

// /api/news/category, /api/news/search,
// /api/news/nearby

func NewsServiceController(appGroup *gin.RouterGroup, db *gorm.DB, redis *redis.Client) {
	app := appGroup.Group("/news")
	controller := _NewController(db, redis)
	app.POST("", controller.createNewsArticle)
	app.GET("/category", controller.getNewsArticlesByCategory)
	app.GET("/score", controller.getNewsArticlesByScore)
	app.GET("/source", controller.getNewsArticlesBySource)
	app.GET("/nearby", controller.getNearByNewsArticle)
	app.GET("/search", controller.getNewsArticleBySearch)
	// app.GET("/trending", controller.getTrendingNewsNearByMe)

}

func (controller *NewsController) createNewsArticle(ctx *gin.Context) {
	var newsArticleInput []NewsArticleUserInuut
	bindingErro := ctx.ShouldBindJSON(&newsArticleInput)
	if bindingErro != nil {
		response.JSONResponse(ctx, bindingErro, nil, nil, "", 0, "")
		return
	}

	newsArticle := ConvertUserInputToNewsArticle(&newsArticleInput)

	service := _NewService(controller.DB, newsArticle)
	err := service.CreateNewsArticle()
	if err != nil {
		response.JSONResponse(ctx, err, nil, nil, "", 0, "")
	}
	ctx.JSON(200, "succeed")
}

func (controller *NewsController) getNewsArticlesByCategory(ctx *gin.Context) {
	var newsArticle []NewsArticle
	category, _ := ctx.GetQuery("name")

	cacheKey := fmt.Sprintf("v1:news:category:%v", category)
	// check if data is rpesent in cache or not
	val, err := redisservice.GetValueFromRedis(controller.Redis, cacheKey)
	if err == redis.Nil {
		var ans response.FinalResponse
		temp, _ := json.Marshal(val)
		json.Unmarshal(temp, &ans)
		ctx.JSON(200, ans)
		return
	}
	fmt.Println("db me jaa rha hain...")
	service := _NewService(controller.DB, &newsArticle)
	result, totalRecord, queryDetail, err := service.GetNewsArticlesByCategory(category)
	if err != nil {
		response.JSONResponse(ctx, err, result, nil, "", 0, "")
		return
	}
	ConvertResultInProperFormatAndReturn(ctx, result, controller.Redis, cacheKey, totalRecord, queryDetail)

	// response.JSONResponse(ctx, err, Convert_NewsArticle_To_NewsArticleResponse(result), controller.Redis, cacheKey)
}

func (controller *NewsController) getNewsArticlesByScore(ctx *gin.Context) {
	var newsArticle []NewsArticle
	x, founded := ctx.GetQuery("val")
	if !founded {
		response.JSONResponse(ctx, fmt.Errorf("query param score is not present"), nil, nil, "", 0, "")
		return
	}
	score, err := strconv.ParseFloat(x, 64)
	if err != nil {
		response.JSONResponse(ctx, err, nil, nil, "", 0, "")
		return
	}

	cacheKey := fmt.Sprintf("v1:news:score:%v", score)
	// check if data is rpesent in cache or not
	val, err := redisservice.GetValueFromRedis(controller.Redis, cacheKey)
	if err == redis.Nil {
		var ans response.FinalResponse
		temp, _ := json.Marshal(val)
		json.Unmarshal(temp, &ans)
		ctx.JSON(200, ans)
		return
	}
	service := _NewService(controller.DB, &newsArticle)
	result, totalRecord, queryDetail, err := service.GetNewsArticlesByScore(score)
	if err != nil {
		response.JSONResponse(ctx, err, result, nil, "", 0, "")
		return
	}
	// response.JSONResponse(ctx, err, Convert_NewsArticle_To_NewsArticleResponse(result))
	ConvertResultInProperFormatAndReturn(ctx, result, controller.Redis, cacheKey, totalRecord, queryDetail)
}

func (controller *NewsController) getNewsArticlesBySource(ctx *gin.Context) {
	var newsArticle []NewsArticle
	source, _ := ctx.GetQuery("val")
	// source, _ := strconv.ParseFloat(x, 64)
	fmt.Println("source ", source)
	// source = "ANI News"

	cacheKey := fmt.Sprintf("v1:news:source:%v", source)
	// check if data is rpesent in cache or not
	val, err := redisservice.GetValueFromRedis(controller.Redis, cacheKey)
	if err == redis.Nil {
		var ans response.FinalResponse
		temp, _ := json.Marshal(val)
		json.Unmarshal(temp, &ans)
		ctx.JSON(200, ans)
		return
	}
	service := _NewService(controller.DB, &newsArticle)
	result, totalRecord, queryDetail, err := service.GetNewsArticlesBySource(source)
	if err != nil {
		response.JSONResponse(ctx, err, result, nil, "", 0, "")
		return
	}
	// response.JSONResponse(ctx, err, Convert_NewsArticle_To_NewsArticleResponse(result))
	ConvertResultInProperFormatAndReturn(ctx, result, controller.Redis, cacheKey, totalRecord, queryDetail)
}

func (controller *NewsController) getNearByNewsArticle(ctx *gin.Context) {
	var newsArticle []NewsArticle
	x, founded := ctx.GetQuery("lat")
	if !founded {
		response.JSONResponse(ctx, fmt.Errorf("query param lat is not present"), nil, nil, "", 0, "")
		return
	}
	y, founded := ctx.GetQuery("lon")
	if !founded {
		response.JSONResponse(ctx, fmt.Errorf("query param lon is not present"), nil, nil, "", 0, "")
		return
	}
	z, founded := ctx.GetQuery("radius")
	if !founded {
		response.JSONResponse(ctx, fmt.Errorf("query param radius is not present"), nil, nil, "", 0, "")
		return
	}
	lat, err := strconv.ParseFloat(x, 64)
	if err != nil {
		response.JSONResponse(ctx, err, nil, nil, "", 0, "")
	}
	lon, err := strconv.ParseFloat(y, 64)
	if err != nil {
		response.JSONResponse(ctx, err, nil, nil, "", 0, "")
	}
	radius, err := strconv.ParseFloat(z, 64)
	if err != nil {
		response.JSONResponse(ctx, err, nil, nil, "", 0, "")
	}
	fmt.Println("lat, log, radius", lat, lon, radius)

	cacheKey := fmt.Sprintf("v1:news:nearby:lat%v:lon%v:radius%v", lat, lon, radius)
	// check if data is rpesent in cache or not
	val, err := redisservice.GetValueFromRedis(controller.Redis, cacheKey)
	if err == redis.Nil {
		var ans response.FinalResponse
		temp, _ := json.Marshal(val)
		json.Unmarshal(temp, &ans)
		ctx.JSON(200, ans)
		return
	}
	service := _NewService(controller.DB, &newsArticle)
	result, totalRecord, queryDetail, err := service.GetNearByNewsArticle(lat, lon, radius)
	if err != nil {
		response.JSONResponse(ctx, err, result, nil, "", 0, "")
		return
	}
	ConvertResultInProperFormatAndReturn(ctx, result, controller.Redis, cacheKey, totalRecord, queryDetail)
}

func (controller *NewsController) getNewsArticleBySearch(ctx *gin.Context) {
	query, founded := ctx.GetQuery("query")
	if !founded {
		response.JSONResponse(ctx, fmt.Errorf("query param search not founded"), nil, nil, "", 0, "")
		return
	}

	llmOutput, err := llmservice.FindLLMEntity(query)
	if err != nil {
		response.JSONResponse(ctx, err, nil, nil, "", 0, "")
		return
	}
	cacheKey := fmt.Sprintf("v1:news:query:%v", llmOutput.Entities)
	// check if data is rpesent in cache or not
	val, err := redisservice.GetValueFromRedis(controller.Redis, cacheKey)
	if err == redis.Nil {
		var ans response.FinalResponse
		temp, _ := json.Marshal(val)
		json.Unmarshal(temp, &ans)
		ctx.JSON(200, ans)
		return
	}
	fmt.Println("db me jaa rha h bhai")
	service := _NewService(controller.DB, &[]NewsArticle{})
	res, totalRecord, queryDetail, err := service.GetNewsArticleBySearch(llmOutput)
	if err != nil {
		response.JSONResponse(ctx, err, nil, nil, "", 0, "")
		return
	}
	ConvertResultInProperFormatAndReturn(ctx, res, controller.Redis, cacheKey, totalRecord, queryDetail)
}

func ConvertResultInProperFormatAndReturn(ctx *gin.Context, result *[]NewsArticle, redisClient *redis.Client, cacheKey string, totalRecords int64, queryDetails string) {
	finalResult := Convert_NewsArticle_To_NewsArticleResponse(result)
	// generate llm summery
	// for i := 0; i < len(*finalResult); i++ {
	// 	summery := llmservice.GenerateSummeryLLM((*finalResult)[i].Title, (*finalResult)[i].Description)
	// 	(*finalResult)[i].LLMSummery = summery
	// }
	response.JSONResponse(ctx, nil, finalResult, redisClient, cacheKey, totalRecords, queryDetails)
}
