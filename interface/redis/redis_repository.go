package redis

import (
	"encoding/json"
	"github.com/drill-shishamo-alliance/asotech_server/model/db"
	"github.com/go-redis/redis"
)

type redisRepository struct {
	*redis.Client
}

type IRedisRepository interface {
	CreateUserRoom(userId, value string) (string, error)
	CreateRoomMember(roomId, userId string) error
	CreateRoomMemberLimit(roomId, num string) error
	GetRoomMemberLimit(roomId string) (*db.RoomMemberLimit, error)
	CreateDemon(roomId, userId string) error
	CreateUserLocation(userId string) error
	DeleteUserLocation(userId string) error
	GetRoomMember(roomId string) (*db.RoomUser, error)
	GetDemonId(roomId string) (string, error)
	GetUserLocation(userId string) (*db.UserLocation, error)
	UpdateRoomMember(roomId string, value []byte) error
}

func NewRedisRepository(kvs *redis.Client) IRedisRepository {
	return &redisRepository{kvs}
}

func (r *redisRepository) CreateUserRoom(userId, value string) (string, error) {
	UserRoomKey := userId + "_ROOM"
	UserRoomValue := value
	// ルームキー
	err := r.Client.Set(UserRoomKey, UserRoomValue, 0).Err()
	if err != nil {
		return "", err
	}
	return UserRoomValue, nil
}

func (r *redisRepository) CreateRoomMember(roomId, userId string) error {
	UserRoomMemberKey := roomId + "_MEMBER"
	// 2. 値を更新
	roomUsers := db.RoomUser{}
	roomUsers.UserId = append(roomUsers.UserId, userId)
	result, _ := json.Marshal(roomUsers)
	err := r.Client.Set(UserRoomMemberKey, result, 0).Err()
	if err != nil {
		return err
	}
	return nil
}

func (r *redisRepository) CreateDemon(roomId, userId string) error {
	UserRoomDemonKey := roomId + "_DEMON"
	err := r.Client.Set(UserRoomDemonKey, userId, 0).Err()
	if err != nil {
		return err
	}
	return nil
}

func (r *redisRepository) GetRoomMember(roomId string) (*db.RoomUser, error) {
	UserRoomMemberKey := roomId + "_MEMBER"
	println("ユーザを取得")
	roomUsersValue, err := r.Client.Get(UserRoomMemberKey).Result()
	if err != nil {
		return nil, err
	}
	// 2. 値を更新
	println("ユーザをデコード")
	roomUsers := db.RoomUser{}
	err = json.Unmarshal([]byte(roomUsersValue), &roomUsers)
	if err != nil {
		return nil,err
	}
	println("デコード完了")
	return &roomUsers, nil
}

func (r *redisRepository) UpdateRoomMember(roomId string, value []byte) error {
	UserRoomMemberKey := roomId + "_MEMBER"
	err := r.Client.Set(UserRoomMemberKey, value, 0).Err()
	if err != nil {
		return err
	}
	return nil
}

func (r *redisRepository) CreateUserLocation(userId string) error {
	UserLocationKey := userId + "_LOCATION"
	err := r.Client.Set(UserLocationKey, "", 0).Err()
	if err != nil {
		return err
	}
	return nil
}

func(r *redisRepository) DeleteUserLocation(userId string) error {
	UserLocationKey := userId + "_LOCATION"
	err := r.Client.Del(UserLocationKey).Err()
	if err != nil {
		return err
	}
	return nil
}

func (r redisRepository) GetDemonId(roomId string) (string, error) {
	UserRoomDemonKey := roomId + "_DEMON"
	userRoomDemonValue, err := r.Client.Get(UserRoomDemonKey).Result()
	if err != nil {
		return "", err
	}
	return userRoomDemonValue, nil
}

func (r *redisRepository) GetUserLocation(userId string) (*db.UserLocation, error) {
	UserLocationKey := userId + "_LOCATION"
	userLocationValue, err := r.Client.Get(UserLocationKey).Result()
	if err != nil {
		return nil, err
	}
	userLocation := db.UserLocation{}
	err = json.Unmarshal([]byte(userLocationValue), &userLocation)
	if err != nil {
		return nil,err
	}
	return &userLocation, nil
}

func  (r *redisRepository) CreateRoomMemberLimit(roomId, num string) error {
	RoomMemberLimitKey := roomId + "_ROOM_MEMBER_LIMIT"
	RoomMemberLimitValue := num
	// ルームキー
	err := r.Client.Set(RoomMemberLimitKey, RoomMemberLimitValue, 0).Err()
	if err != nil {
		return err
	}
	return nil
}

func (r *redisRepository) GetRoomMemberLimit(roomId string) (*db.RoomMemberLimit, error) {
	RoomMemberLimitKey := roomId + "_ROOM_MEMBER_LIMIT"
	println(RoomMemberLimitKey)
	roomMemberLimitValue, err := r.Client.Get(RoomMemberLimitKey).Result()
	if err != nil {
		return nil, err
	}
	roomMemberLimit := new(db.RoomMemberLimit)
	roomMemberLimit.Value = roomMemberLimitValue
	return roomMemberLimit, nil
}