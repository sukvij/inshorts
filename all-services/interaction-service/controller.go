package interactionservice

import (
	"github.com/gin-gonic/gin"
	"github.com/sukvij/inshorts/inshortfers/response"
	"gorm.io/gorm"
)

type InteractionController struct {
	DB          *gorm.DB
	Interaction *[]UserInteraction
}

func _NewController(db *gorm.DB) *InteractionController {
	return &InteractionController{DB: db}
}

func UserInteractionController(appGroup *gin.RouterGroup, db *gorm.DB) {
	// /v1/interaction
	app := appGroup.Group("/interaction")
	controller := _NewController(db)
	app.POST("", controller.createUserInteraction)

}

func (controller *InteractionController) createUserInteraction(ctx *gin.Context) {
	var interactions []UserInteraction
	err := ctx.ShouldBindJSON(&interactions)
	if err != nil {
		response.JSONResponse(ctx, err, nil)
		return
	}
	service := _NewService(controller.DB, &interactions)
	err = service.CreateUserInteraction()
	response.JSONResponse(ctx, err, nil)
}
