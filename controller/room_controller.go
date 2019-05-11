package controller

import (
	"github.com/drill-shishamo-alliance/asotech_server/model"
	"github.com/gin-gonic/gin"
	"net/http"
)

type roomController struct {
	model.IRoom
}

type IRoomController interface {
	CreateTheRoom(ctx *gin.Context)
	BelongToTheRoom(ctx *gin.Context)
	IsAllMemberReady(ctx *gin.Context)
	GetRemainingHumans(ctx *gin.Context)
	GetHumansLocation(ctx *gin.Context)
	GetHumanCollaborate(ctx *gin.Context)
	GetDemonsLocation(ctx *gin.Context)
}

func NewRoomController(r model.IRoom) IRoomController {
	return &roomController{r}
}

func (ctrl *roomController) CreateTheRoom(ctx *gin.Context) {
	userId := ctx.GetHeader("user_id")
	result, err := ctrl.IRoom.Insert(userId, "3600", 5)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"status": err})
	}
	ctx.JSON(http.StatusCreated, result)
}

func (ctrl *roomController) IsAllMemberReady(ctx *gin.Context) {
	roomId := ctx.Param("id")
	result, err := ctrl.IRoom.IsMemberReady(roomId)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"status": err})
	}
	ctx.JSON(http.StatusBadRequest, gin.H{"status": result})
}

func (ctrl *roomController) BelongToTheRoom(ctx *gin.Context) {
	roomId := ctx.Param("id")
	userId := ctx.GetHeader("user_id")
	err := ctrl.IRoom.UpdateMember(roomId, userId)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"status": err})
	}
	ctx.JSON(http.StatusOK, gin.H{})
}

func (ctrl *roomController) GetRemainingHumans(ctx *gin.Context) {
	roomId := ctx.Param("id")
	result, err := ctrl.IRoom.SelectRemainingHuman(roomId)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"status": err})
	}
	ctx.JSON(http.StatusOK, result)
}

func (ctrl *roomController) GetHumanCollaborate(ctx *gin.Context) {
	roomId := ctx.Param("id")
	userId := ctx.GetHeader("user_id")
	result, err := ctrl.IRoom.SelectCollaborateHuman(roomId, userId, 80)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"status": err})
	}
	ctx.JSON(http.StatusOK, result)
}

func (ctrl *roomController) GetHumansLocation(ctx *gin.Context) {
	roomId := ctx.Param("id")
	result, err := ctrl.IRoom.SelectHumansLocation(roomId)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"status": err})
	}
	ctx.JSON(http.StatusOK, result)
}

func (ctrl *roomController) GetDemonsLocation(ctx *gin.Context) {
	roomId := ctx.Param("id")
	result, err := ctrl.IRoom.SelectRoomMember(roomId)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"status": err})
	}
	ctx.JSON(http.StatusOK, result)
}