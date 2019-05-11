package controller

import "github.com/gin-gonic/gin"

type roomController struct {}

type IRoomController interface {
	Generate(ctx *gin.Context)
	IsRoomMemberAlready(ctx *gin.Context)
	Join(ctx *gin.Context)
	GetRemainingHumans(ctx *gin.Context)
	GetHumansLocation(ctx *gin.Context)
	GetDemonsLocation(ctx *gin.Context)
}

func NewRoomController() IRoomController {
	return &roomController{}
}

func (ctrl *roomController) Generate(ctx *gin.Context) {

}