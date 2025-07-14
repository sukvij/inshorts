package allservices

import (
	"github.com/gin-gonic/gin"
	interactionservice "github.com/sukvij/inshorts/all-services/interaction-service"
	newsservice "github.com/sukvij/inshorts/all-services/news-service"
	"github.com/sukvij/inshorts/inshortfers/logs"
	"go.opentelemetry.io/otel/sdk/trace"
	"gorm.io/gorm"
)

func RouteService(engine *gin.Engine, db *gorm.DB, logs *logs.AgreeGateLoager, tracker *trace.TracerProvider) {
	app := engine.Group("/v1")
	newsservice.NewsServiceController(app, db) //        /v1/news-article   --> this is the endpoints
	interactionservice.UserInteractionController(app, db)
}
