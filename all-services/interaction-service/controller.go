package interactionservice

import (
	"fmt"

	"github.com/gin-gonic/gin"
	"github.com/go-redis/redis"
	llmservice "github.com/sukvij/inshorts/all-services/llm-service"
	newsservice "github.com/sukvij/inshorts/all-services/news-service"
	"github.com/sukvij/inshorts/inshortfers/logs"
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
		response.JSONResponse(ctx, err, nil)
		return
	}
	service := _NewService(controller.DB, &interactions, controller.Redis)
	err = service.CreateUserInteraction()
	if err != nil {
		controller.Logs.Error(err)
	}
	response.JSONResponse(ctx, err, nil)
}

func (controller *InteractionController) trendingNewsArticles(ctx *gin.Context) {
	lat := ctx.Query("lat")
	lon := ctx.Query("lon")
	limit := ctx.Query("limit")
	fmt.Println(lat, lon, limit)
	service := _NewService(controller.DB, &[]UserInteraction{}, controller.Redis)
	result, err := service.TrendingNewsArticles(lat, lon, limit)
	if err != nil {
		controller.Logs.Error(err)
		response.JSONResponse(ctx, err, nil)
		return
	}
	ConvertResultInProperFormatAndReturn(ctx, result)
}

func ConvertResultInProperFormatAndReturn(ctx *gin.Context, result *[]newsservice.NewsArticle) {
	finalResult := newsservice.Convert_NewsArticle_To_NewsArticleResponse(result)
	// generate llm summery
	for i := 0; i < len(*finalResult); i++ {
		summery := llmservice.GenerateSummeryLLM((*finalResult)[i].Title, (*finalResult)[i].Description)
		(*finalResult)[i].LLMSummery = summery
	}
	response.JSONResponse(ctx, nil, finalResult)
}
