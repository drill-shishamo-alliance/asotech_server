package router

import (
	"github.com/drill-shishamo-alliance/asotech_server/interface/di"
	"github.com/gin-gonic/gin"
)

type routerHandler struct {}

type IRouterHandler interface {
	SetUpRouter() *gin.Engine
}

func NewRouterHandler() IRouterHandler {
	return &routerHandler{}
}

func (r *routerHandler) SetUpRouter() *gin.Engine {
	/**
	 * Set Up Dependency Injection
	 */
	ctrl := di.InitRoomController()
	/*
	 * SetUp Routing
	 */
	router := gin.Default()
	router.POST("/rooms", ctrl.CreateTheRoom)
	router.GET("rooms/{id}/checkin/status", ctrl.IsAllMemberReady)
	router.POST("rooms/{id}/checkin", ctrl.BelongToTheRoom)
	router.GET("rooms/{id}/remaining/human", ctrl.GetRemainingHumans)
	router.GET("rooms/{id}/humans/collaborate", ctrl.GetHumanCollaborate)
	router.GET("rooms/{id}/humans/locations", ctrl.GetHumansLocation)
	router.GET("rooms/{id}/demons/locations", ctrl.GetDemonsLocation)

	return router
}