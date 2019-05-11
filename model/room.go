package model

import (
	"encoding/json"
	"github.com/drill-shishamo-alliance/asotech_server/interface/guid"
	"github.com/drill-shishamo-alliance/asotech_server/interface/redis"
	"github.com/drill-shishamo-alliance/asotech_server/model/db"
)

type room struct {
	redis.IRedisRepository
	guid.IGuidUtil
}

type IRoom interface {
	Insert(userId, restrictTime string, memberNum int) (string, error)
	UpdateMember(roomId, userId string) (string, error)
	SelectRoomMember(roomId string) (string, error)
	DecreaseMemberStatus(roomId, userId string) error
	SelectRemainingHuman(roomId string) (string, error)
	SelectHumansLocation(roomId string) ([]*db.UserLocation, error)
	SelectDemonsLocation(roomId string) (*db.UserLocation, error)
}

func NewRoom(repo redis.IRedisRepository , id guid.IGuidUtil) IRoom {
	return &room{repo, id}
}

// Insert
func (r *room) Insert(userId, restrictTime string, memberNum int) (string, error) {
	// ゲームの部屋作成
	UserRoomValue, err := r.IRedisRepository.CreateUserRoom(userId, r.IGuidUtil.CreateGuid())
	if err != nil {
		return "", nil
	}
	// ゲームのメンバー
	err =  r.IRedisRepository.CreateRoomMember(UserRoomValue)
	if err != nil {
		return "", nil
	}
	// 鬼の決定
	err = r.IRedisRepository.CreateDemon(UserRoomValue, userId)
	if err != nil {
		return "", nil
	}
	return UserRoomValue, nil
}

func (r *room) UpdateMember(roomId, userId string) (string, error) {
	// 部屋にユーザを追加
	// 1. ユーザのデータを取得
	roomMember, err := r.IRedisRepository.GetRoomMember(roomId)
	if err != nil {
		return "", err
	}
	// 2. データ追加
	roomMember.UserId = append(roomMember.UserId, userId)
	roomMemberNewValue, err := json.Marshal(roomMember)
	if err != nil {
		return "", err
	}
	// 3. 再保存
	err = r.IRedisRepository.UpdateRoomMember(roomId, roomMemberNewValue)
	if err != nil {
		return "", err
	}
	// 4. ユーザの位置情報保存
	err = r.IRedisRepository.CreateUserLocation(userId)
	if err != nil {
		return "", err
	}
	return "", nil
}

func (r *room) SelectRoomMember(roomId string) (string, error) {
	return "", nil
}

func (r *room) DecreaseMemberStatus(roomId, userId string) error {
	// 1. メンバーを減らす
	roomMember, err := r.IRedisRepository.GetRoomMember(roomId)
	if err != nil {
		return err
	}
	// 2. 値を更新
	result := db.RoomUser{}
	for _, values := range roomMember.UserId {
		if values != userId {
			result.UserId = append(result.UserId, values)
		}
	}
	// 3. 値を保存
	roomMemberNewValue, err := json.Marshal(result)
	err = r.IRedisRepository.UpdateRoomMember(roomId, roomMemberNewValue)
	if err != nil {
		return err
	}
	// 4. ロケーション消去
	err = r.IRedisRepository.DeleteUserLocation(userId)
	if err != nil {
		return err
	}
	return nil
}

func (r *room) SelectRemainingHuman(roomId string) (string, error) {
	// 1. ルームのメンバーを取得
	roomMemberValue, err := r.IRedisRepository.GetRoomMember(roomId)
	if err != nil {
		return "", err
	}
	// 2-1. 鬼のIDを取得
	demonId, err := r.IRedisRepository.GetDemonId(roomId)
	if err != nil {
		return "", err
	}
	count := 0
	for _, values := range roomMemberValue.UserId {
		if values != demonId {
			count ++
		}
	}
	return string(count), nil
}

func (r *room) SelectHumansLocation(roomId string) ([]*db.UserLocation, error) {
	// 1. ルームのメンバーを取得
	roomMemberValue, err := r.IRedisRepository.GetRoomMember(roomId)
	if err != nil {
		return nil, err
	}
	// 2-1. 鬼のIDを取得
	demonId, err := r.IRedisRepository.GetDemonId(roomId)
	if err != nil {
		return nil, err
	}
	// 3. 位置情報を取得
	var result []*db.UserLocation
	for _, value := range roomMemberValue.UserId {
		if value != demonId {
			userLocation, err := r.IRedisRepository.GetUserLocation(value)
			if err != nil {
				return nil, err
			}
			result = append(result, userLocation)
		}
	}
	return result, nil
}

func (r *room) SelectDemonsLocation(roomId string) (*db.UserLocation, error) {
	// 1. ルームのメンバーを取得
	roomMemberValue, err := r.IRedisRepository.GetRoomMember(roomId)
	if err != nil {
		return nil, err
	}
	// 2-1. 鬼のIDを取得
	demonId, err := r.IRedisRepository.GetDemonId(roomId)
	if err != nil {
		return nil, err
	}
	// 3. 位置情報を取得
	result := db.UserLocation{}
	for _, value := range roomMemberValue.UserId {
		if value == demonId {
			userLocation, err := r.IRedisRepository.GetUserLocation(value)
			if err != nil {
				return nil, err
			}
			result = *userLocation
			break
		}
	}
	return &result, nil
}
