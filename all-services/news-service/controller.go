package newsservice

import (
	"encoding/json"

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

func NewsServiceController(appGroup *gin.RouterGroup, db *gorm.DB) {
	app := appGroup.Group("/news-article")
	controller := _NewController(db)
	app.POST("", controller.createNewsArticle)
	app.GET("", controller.getNewsArticles)
}

func (controller *NewsController) createNewsArticle(ctx *gin.Context) {
	var newsArticleInput []NewsArticleUserInuut
	bindingErro := ctx.ShouldBindJSON(&newsArticleInput)
	if bindingErro != nil {
		ctx.JSON(200, "binding erro")
		return
	}

	var newsArticle []NewsArticle
	for _, input := range newsArticleInput {
		var temp NewsArticle
		res, _ := json.Marshal(input)
		json.Unmarshal(res, &temp)
		x, _ := json.Marshal((input.Category))
		temp.Category = x
		// fmt.Println("Cat", string(temp.Category))
		newsArticle = append(newsArticle, temp)
	}
	service := _NewService(controller.DB, &newsArticle)
	errs := service.CreateNewsArticle()
	ctx.JSON(200, errs)
}

func (controller *NewsController) getNewsArticles(ctx *gin.Context) {
	var newsArticleInput []NewsArticleUserInuut
	bindingErro := ctx.ShouldBindJSON(&newsArticleInput)
	if bindingErro != nil {
		ctx.JSON(200, "binding erro")
		return
	}
	// service := _NewService(controller.DB, &newsArticleInput)
	// result, _ := service.GetNewsArticles()
	// ctx.JSON(200, newsArticle)
}
