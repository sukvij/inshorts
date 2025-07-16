package newsservice

import (
	"fmt"
	"strconv"

	"github.com/gin-gonic/gin"
	llmservice "github.com/sukvij/inshorts/all-services/llm-service"
	"github.com/sukvij/inshorts/inshortfers/response"
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
	// app.GET("/trending", controller.getTrendingNewsNearByMe)

}

func (controller *NewsController) createNewsArticle(ctx *gin.Context) {
	var newsArticleInput []NewsArticleUserInuut
	bindingErro := ctx.ShouldBindJSON(&newsArticleInput)
	if bindingErro != nil {
		response.JSONResponse(ctx, bindingErro, nil)
		return
	}

	newsArticle := ConvertUserInputToNewsArticle(&newsArticleInput)

	service := _NewService(controller.DB, newsArticle)
	err := service.CreateNewsArticle()
	response.JSONResponse(ctx, err, nil)
}

func (controller *NewsController) getNewsArticlesByCategory(ctx *gin.Context) {
	var newsArticle []NewsArticle
	category, _ := ctx.GetQuery("name")
	fmt.Println(category)
	service := _NewService(controller.DB, &newsArticle)
	result, err := service.GetNewsArticlesByCategory(category)
	if err != nil {
		response.JSONResponse(ctx, err, result)
		return
	}
	response.JSONResponse(ctx, err, Convert_NewsArticle_To_NewsArticleResponse(result))
}

func (controller *NewsController) getNewsArticlesByScore(ctx *gin.Context) {
	var newsArticle []NewsArticle
	x, founded := ctx.GetQuery("val")
	if !founded {
		response.JSONResponse(ctx, fmt.Errorf("query param score is not present"), nil)
		return
	}
	score, err := strconv.ParseFloat(x, 64)
	if err != nil {
		response.JSONResponse(ctx, err, nil)
		return
	}
	service := _NewService(controller.DB, &newsArticle)
	result, err := service.GetNewsArticlesByScore(score)
	if err != nil {
		response.JSONResponse(ctx, err, result)
		return
	}
	response.JSONResponse(ctx, err, Convert_NewsArticle_To_NewsArticleResponse(result))
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
		response.JSONResponse(ctx, err, result)
		return
	}
	response.JSONResponse(ctx, err, Convert_NewsArticle_To_NewsArticleResponse(result))
}

func (controller *NewsController) getNearByNewsArticle(ctx *gin.Context) {
	var newsArticle []NewsArticle
	x, founded := ctx.GetQuery("lat")
	if !founded {
		response.JSONResponse(ctx, fmt.Errorf("query param lat is not present"), nil)
		return
	}
	y, founded := ctx.GetQuery("lon")
	if !founded {
		response.JSONResponse(ctx, fmt.Errorf("query param lon is not present"), nil)
		return
	}
	z, founded := ctx.GetQuery("radius")
	if !founded {
		response.JSONResponse(ctx, fmt.Errorf("query param radius is not present"), nil)
		return
	}
	lat, err := strconv.ParseFloat(x, 64)
	if err != nil {
		response.JSONResponse(ctx, err, nil)
	}
	lon, err := strconv.ParseFloat(y, 64)
	if err != nil {
		response.JSONResponse(ctx, err, nil)
	}
	radius, err := strconv.ParseFloat(z, 64)
	if err != nil {
		response.JSONResponse(ctx, err, nil)
	}
	fmt.Println("lat, log, radius", lat, lon, radius)

	service := _NewService(controller.DB, &newsArticle)
	result, err := service.GetNearByNewsArticle(lat, lon, radius)
	if err != nil {
		response.JSONResponse(ctx, err, result)
		return
	}
	updatedResponse := Convert_NewsArticle_To_NewsArticleResponse(result)
	// var summaries []string
	// for i := 0; i < len(*updatedResponse); i++ {
	// 	haha := llmservice.GenerateSummeryLLM((*updatedResponse)[i].Title, (*updatedResponse)[i].Description)
	// 	(*updatedResponse)[i].LLMSummery = haha
	// }
	response.JSONResponse(ctx, err, updatedResponse)
}

func (controller *NewsController) getNewsArticleBySearch(ctx *gin.Context) {
	query, founded := ctx.GetQuery("query")
	if !founded {
		response.JSONResponse(ctx, fmt.Errorf("query param search not founded"), nil)
		return
	}

	llmOutput, err := llmservice.FindLLMEntity(query)
	if err != nil {
		response.JSONResponse(ctx, err, nil)
		return
	}
	service := _NewService(controller.DB, &[]NewsArticle{})
	res, err := service.GetNewsArticleBySearch(llmOutput)
	if err != nil {
		response.JSONResponse(ctx, err, nil)
		return
	}
	response.JSONResponse(ctx, err, Convert_NewsArticle_To_NewsArticleResponse(res))
	// ConvertResultInProperFormatAndReturn(ctx, re)
}

func ConvertResultInProperFormatAndReturn(ctx *gin.Context, result *[]NewsArticle) {
	finalResult := Convert_NewsArticle_To_NewsArticleResponse(result)
	// generate llm summery
	for i := 0; i < len(*finalResult); i++ {
		summery := llmservice.GenerateSummeryLLM((*finalResult)[i].Title, (*finalResult)[i].Description)
		(*finalResult)[i].LLMSummery = summery
	}
	response.JSONResponse(ctx, nil, finalResult)
}
