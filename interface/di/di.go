package di

import (
	"github.com/drill-shishamo-alliance/asotech_server/controller"
	"github.com/drill-shishamo-alliance/asotech_server/interface/guid"
	"github.com/drill-shishamo-alliance/asotech_server/interface/redis"
	"github.com/drill-shishamo-alliance/asotech_server/model"
)

func InitRoomController() controller.IRoomController {
	redisProvider := redis.NewRedisProvider()
	redisRepo := redis.NewRedisRepository(redisProvider.CreateClient())
	roomModel := model.NewRoom(redisRepo, guid.NewGuidUtil())
	return controller.NewRoomController(roomModel)
}