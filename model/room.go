package model

import (
	"encoding/json"
	"github.com/drill-shishamo-alliance/asotech_server/interface/guid"
	"github.com/drill-shishamo-alliance/asotech_server/model/db"
	"github.com/go-redis/redis"
	"strconv"
)

type room struct {
	*redis.Client
	guid.IGuidUtil
}

type IRoom interface {
	Insert(userId string, restrictTime int)
	UpdateMember(roomId, userId string) (string, error)
	SelectRoomMember(roomId string) (string, error)
	UpdateMemberStatus(userId string) (string, error)
	SelectRemainingHuman(roomId string) (string, error)
	SelectHumansLocation(roomId string) (string, error)
	SelectDemonsLocation(roomId string) (string, error)
}

func NewRoom(kvs *redis.Client, id guid.IGuidUtil) IRoom {
	return &room{kvs, id}
}

// Insert
func (r *room) Insert(userId, restrictTime string, memberNum int) (string, error) {
	UserRoomKey := userId + "_ROOM"
	UserRoomValue := r.IGuidUtil.CreateGuid()
	// ルームキー
	err := r.Client.Set(UserRoomKey, UserRoomValue, 0).Err()
	if err != nil {
		return "", err
	}
	// ゲームのメンバー
	UserRoomMemberKey := UserRoomValue + "_MEMBER"
	err = r.Client.Set(UserRoomMemberKey, "", 0).Err()
	if err != nil {
		return "", err
	}
	// 鬼の決定
	UserRoomDemonKey := UserRoomValue + "_DEMON"
	err = r.Client.Set(UserRoomDemonKey, userId, 0).Err()
	if err != nil {
		return "", err
	}
	return UserRoomValue, nil
}

func (r *room) UpdateMember(roomId, userId string) (string, error) {
	// 部屋にユーザを追加
	// 1. ユーザのデータを取得
	UserRoomMemberKey := roomId + "_MEMBER"
	roomUsersValue, err := r.Client.Get(UserRoomMemberKey).Result()
	if err != nil {
		return "", err
	}
	// 2. 値を更新
	roomUsers := db.RoomUser{}
	err = json.Unmarshal([]byte(roomUsersValue), &roomUsers)
	if err != nil {
		return "", nil
	}
	roomUsers.UserId = append(roomUsers.UserId, userId)
	roomUsersNewValue, err := json.Marshal(roomUsers)
	if err != nil {
		return "", nil
	}
	// 3. 再保存
	err = r.Client.Set(UserRoomMemberKey, roomUsersNewValue, 0).Err()
	if err != nil {
		return "", err
	}
	// 4. ユーザの位置情報保存
	UserLocationKey := userId + "_LOCATION"
	err = r.Client.Set(UserLocationKey, "", 0).Err()
	if err != nil {
		return "", err
	}
	return "", nil
}

func (r *room) SelectRoomMember(roomId string) (string, error){
	// 1. ユーザのデータを取得
	UserRoomMemberKey := roomId + "_ROOM"
	roomUsersValue, err := r.Client.Get(UserRoomMemberKey).Result()
	if err != nil {
		return "", err
	}
	// 2. 値を更新
	roomUsers := db.RoomUser{}
	err = json.Unmarshal([]byte(roomUsersValue), &roomUsers)
	if err != nil {
		return "", nil
	}
	return string(len(roomUsers.UserId)), nil
}

func (r *room) UpdateMemberStatus(roomId, userId string) error {
	// 1. メンバーを減らす
	UserRoomMemberKey := roomId + "_MEMBER"
	roomUsersValue, err := r.Client.Get(UserRoomMemberKey).Result()
	if err != nil {
		return err
	}
	// 2. 値を更新
	roomUsers := db.RoomUser{}
	err = json.Unmarshal([]byte(roomUsersValue), &roomUsers)
	if err != nil {
		return nil
	}
	result := db.RoomUser{}
	for _, values := range roomUsers.UserId {
		if values != userId {
			result.UserId = append(result.UserId, values)
		}
	}
	// 3. 再保存
	err = r.Client.Set(UserRoomMemberKey, result, 0).Err()
	if err != nil {
		return err
	}
	// 4. ロケーション消去
	UserLocationKey := userId + "_LOCATION"
	err = r.Client.Del(UserLocationKey).Err()
	if err != nil {
		return err
	}
	return nil
}

func (r *room) SelectRemainingHuman(roomId string) (string, error) {
	// 1. ルームのメンバーを取得
	UserRoomMemberKey := roomId + "_MEMBER"
	roomUsersValue, err := r.Client.Get(UserRoomMemberKey).Result()
	if err != nil {
		return "", err
	}
	// 2. 鬼以外の生き残りを取得
	roomUsers := db.RoomUser{}
	err = json.Unmarshal([]byte(roomUsersValue), &roomUsers)
	if err != nil {
		return "", nil
	}
	// 2-1. 鬼のIDを取得
	UserRoomDemonKey := roomId + "_DEMON"
	userRoomDemonValue, err := r.Client.Get(UserRoomDemonKey).Result()
	if err != nil {
		return "", err
	}
	count := 0
	for _, values := range roomUsers.UserId {
		if values != userRoomDemonValue {
			count ++
		}
	}
	return string(count), nil
}

func (r *room) SelectHumansLocation(roomId string) (string, error) {
	// 1. ルームのメンバーを取得
	UserRoomMemberKey := roomId + "_MEMBER"
	roomUsersValue, err := r.Client.Get(UserRoomMemberKey).Result()
	if err != nil {
		return "", err
	}
	// 2. 鬼以外の生き残りを取得
	roomUsers := db.RoomUser{}
	err = json.Unmarshal([]byte(roomUsersValue), &roomUsers)
	if err != nil {
		return "", nil
	}
	// 2-1. 鬼のIDを取得
	UserRoomDemonKey := roomId + "_DEMON"
	userRoomDemonValue, err := r.Client.Get(UserRoomDemonKey).Result()
	if err != nil {
		return "", err
	}
	// 3. 位置情報を取得
	for _, value := range roomUsers.UserId {
		UserLocationKey := value + "_LOCATION"
		userLocation, err = r.Client.Get(UserLocationKey).Result()
		if err != nil {
			return "", err
		}
	}
	return "", nil
}

func (r *room) SelectDemonsLocation(roomId string) (string, error) {
	return "", nil
}
