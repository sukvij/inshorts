package newsservice

import (
	"fmt"
	"strconv"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

type NewsController struct {
	DB          *gorm.DB
	NewsArticle []*NewsArticle
}

func _NewController(db *gorm.DB) *NewsController {
	return &NewsController{DB: db}
}

// /api/news/category, /api/news/search,
// /api/news/nearby

func NewsServiceController(appGroup *gin.RouterGroup, db *gorm.DB) {
	app := appGroup.Group("/news")
	controller := _NewController(db)
	app.POST("", controller.createNewsArticle)
	app.GET("/category", controller.getNewsArticlesByCategory)
	app.GET("/score", controller.getNewsArticlesByScore)
	app.GET("/source", controller.getNewsArticlesBySource)
	app.GET("/nearby", controller.getNearByNewsArticle)
	app.GET("/search", controller.getNewsArticleBySearch)
}

func (controller *NewsController) createNewsArticle(ctx *gin.Context) {
	var newsArticleInput []NewsArticleUserInuut
	bindingErro := ctx.ShouldBindJSON(&newsArticleInput)
	if bindingErro != nil {
		ctx.JSON(200, "binding erro")
		return
	}

	newsArticle := ConvertUserInputToNewsArticle(&newsArticleInput)

	service := _NewService(controller.DB, newsArticle)
	err := service.CreateNewsArticle()
	ctx.JSON(200, err)
}

func (controller *NewsController) getNewsArticlesByCategory(ctx *gin.Context) {
	var newsArticle []NewsArticle
	category, _ := ctx.GetQuery("name")
	fmt.Println(category)
	service := _NewService(controller.DB, &newsArticle)
	result, err := service.GetNewsArticlesByCategory(category)
	if err != nil {
		ctx.JSON(500, err)
		return
	}
	ctx.JSON(200, Convert_NewsArticle_To_NewsArticleResponse(result))
}

func (controller *NewsController) getNewsArticlesByScore(ctx *gin.Context) {
	var newsArticle []NewsArticle
	x, _ := ctx.GetQuery("val")
	score, _ := strconv.ParseFloat(x, 64)
	fmt.Println(score)
	service := _NewService(controller.DB, &newsArticle)
	result, err := service.GetNewsArticlesByScore(score)
	if err != nil {
		ctx.JSON(500, err)
		return
	}
	ctx.JSON(200, result)
}

func (controller *NewsController) getNewsArticlesBySource(ctx *gin.Context) {
	var newsArticle []NewsArticle
	source, _ := ctx.GetQuery("val")
	// source, _ := strconv.ParseFloat(x, 64)
	fmt.Println("source ", source)
	// source = "ANI News"
	service := _NewService(controller.DB, &newsArticle)
	result, err := service.GetNewsArticlesBySource(source)
	if err != nil {
		ctx.JSON(500, err)
		return
	}
	ctx.JSON(200, Convert_NewsArticle_To_NewsArticleResponse(result))
}

func (controller *NewsController) getNearByNewsArticle(ctx *gin.Context) {
	var newsArticle []NewsArticle
	x, _ := ctx.GetQuery("lat")
	y, _ := ctx.GetQuery("lon")
	z, _ := ctx.GetQuery("radius")
	lat, _ := strconv.ParseFloat(x, 64)
	lon, _ := strconv.ParseFloat(y, 64)
	radius, _ := strconv.ParseFloat(z, 64)
	fmt.Println("lat, log, radius", lat, lon, radius)

	service := _NewService(controller.DB, &newsArticle)
	result, err := service.GetNearByNewsArticle(lat, lon, radius)
	if err != nil {
		ctx.JSON(500, err)
		return
	}
	ctx.JSON(200, Convert_NewsArticle_To_NewsArticleResponse(result))
}

func (controller *NewsController) getNewsArticleBySearch(ctx *gin.Context) {

}
