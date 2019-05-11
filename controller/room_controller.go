package controller

import (
	"github.com/drill-shishamo-alliance/asotech_server/model"
	"github.com/drill-shishamo-alliance/asotech_server/model/view"
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
	var form view.PostRoom
	err := ctx.Bind(&form)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"status": err})
	}
	result, err := ctrl.IRoom.Insert(form.UserId, "3600", 5)
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
	var body view.PostBelongToTheRoom
	err := ctx.Bind(&body)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"status": err})
	}
	err = ctrl.IRoom.UpdateMember(roomId, body.UserId)
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
	var body view.PostHumanCollaborate
	err := ctx.Bind(&body)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"status": err})
	}
	result, err := ctrl.IRoom.SelectCollaborateHuman(roomId, body.UserId, 80)
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