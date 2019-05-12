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
	// /users -> apiary
	router.POST("/rooms", ctrl.CreateTheRoom)
	router.GET("/rooms/checkin/status", ctrl.IsAllMemberReady)
	router.POST("/rooms/checkin", ctrl.BelongToTheRoom)
	router.GET("/rooms/remaining/human", ctrl.GetRemainingHumans)
	router.GET("/rooms/humans/collaborate", ctrl.GetHumanCollaborate)
	router.GET("/rooms/humans/locations", ctrl.GetHumansLocation)
	router.GET("/rooms/demons/locations", ctrl.GetDemonsLocation)

	return router
}