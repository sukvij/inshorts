package interactionservice

import (
	"fmt"

	"github.com/gin-gonic/gin"
	"github.com/go-redis/redis"
	newsservice "github.com/sukvij/inshorts/all-services/news-service"
	"github.com/sukvij/inshorts/inshortfers/response"
	"gorm.io/gorm"
)

type InteractionController struct {
	DB          *gorm.DB
	Interaction *[]UserInteraction
	Redis       *redis.Client
}

func _NewController(db *gorm.DB, redis *redis.Client) *InteractionController {
	return &InteractionController{DB: db, Redis: redis}
}

func UserInteractionController(appGroup *gin.RouterGroup, db *gorm.DB, redis *redis.Client) {
	// /v1/interaction
	app := appGroup.Group("/interaction")
	controller := _NewController(db, redis)
	app.POST("", controller.createUserInteraction)
	app.GET("/trending", controller.trendingNewsArticles)
}

func (controller *InteractionController) createUserInteraction(ctx *gin.Context) {
	var interactions []UserInteraction
	err := ctx.ShouldBindJSON(&interactions)
	if err != nil {
		response.JSONResponse(ctx, err, nil)
		return
	}
	service := _NewService(controller.DB, &interactions, controller.Redis)
	err = service.CreateUserInteraction()
	response.JSONResponse(ctx, err, nil)
}

func (controller *InteractionController) trendingNewsArticles(ctx *gin.Context) {
	lat := ctx.Query("lat")
	lon := ctx.Query("lon")
	limit := ctx.Query("limit")
	fmt.Println(lat, lon, limit)
	service := _NewService(controller.DB, &[]UserInteraction{}, controller.Redis)
	result, err := service.TrendingNewsArticles(lat, lon, limit)
	finalResult := newsservice.Convert_NewsArticle_To_NewsArticleResponse(result)
	response.JSONResponse(ctx, err, finalResult)
}
