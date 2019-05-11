package model

import (
	"encoding/json"
	"fmt"
	"github.com/drill-shishamo-alliance/asotech_server/interface/guid"
	"github.com/drill-shishamo-alliance/asotech_server/interface/redis"
	"github.com/drill-shishamo-alliance/asotech_server/model/db"
	"math"
)

type room struct {
	redis.IRedisRepository
	guid.IGuidUtil
}

type IRoom interface {
	Insert(userId, restrictTime string, memberNum int) (*db.RoomId, error)
	UpdateMember(roomId, userId string) error
	SelectRoomMember(roomId string) (string, error)
	DecreaseMemberStatus(roomId, userId string) error
	SelectRemainingHuman(roomId string) (*db.RemainingHuman, error)
	SelectHumansLocation(roomId string) ([]*db.UserLocation, error)
	SelectDemonsLocation(roomId string) (*db.UserLocation, error)
	SelectCollaborateHuman(roomId, userId string, circle float64) (*db.RoomUser ,error)
	IsMemberReady(roomId string) (bool, error)
}

func NewRoom(repo redis.IRedisRepository , id guid.IGuidUtil) IRoom {
	return &room{repo, id}
}

// Insert
func (r *room) Insert(userId, restrictTime string, memberNum int) (*db.RoomId, error) {
	// ゲームの部屋作成
	UserRoomValue, err := r.IRedisRepository.CreateUserRoom(userId, r.IGuidUtil.CreateGuid())
	if err != nil {
		return nil, err
	}
	// ゲームのメンバー
	err =  r.IRedisRepository.CreateRoomMember(UserRoomValue)
	if err != nil {
		return nil, err
	}
	// 鬼の決定
	err = r.IRedisRepository.CreateDemon(UserRoomValue, userId)
	if err != nil {
		return nil, err
	}
	result := new(db.RoomId)
	result.Value = UserRoomValue
	return result, nil
}

func (r *room) UpdateMember(roomId, userId string) error {
	// 部屋にユーザを追加
	// 1. ユーザのデータを取得
	roomMember, err := r.IRedisRepository.GetRoomMember(roomId)
	if err != nil {
		return err
	}
	// 2. データ追加
	roomMember.UserId = append(roomMember.UserId, userId)
	roomMemberNewValue, err := json.Marshal(roomMember)
	if err != nil {
		return err
	}
	// 3. 再保存
	err = r.IRedisRepository.UpdateRoomMember(roomId, roomMemberNewValue)
	if err != nil {
		return err
	}
	// 4. ユーザの位置情報保存
	err = r.IRedisRepository.CreateUserLocation(userId)
	if err != nil {
		return err
	}
	return nil
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

func (r *room) SelectRemainingHuman(roomId string) (*db.RemainingHuman, error) {
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
	count := 0
	for _, values := range roomMemberValue.UserId {
		if values != demonId {
			count ++
		}
	}
	result := new(db.RemainingHuman)
	result.Value = string(count)
	return result, nil
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

func (r *room) SelectCollaborateHuman(roomId, userId string, circle float64) (*db.RoomUser, error) {
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
	myLocation, err := r.IRedisRepository.GetUserLocation(userId)
	if err != nil {
		return nil, err
	}
	// 3. 位置情報を取得
	result := new(db.RoomUser)
	for _, value := range roomMemberValue.UserId {
		if value != demonId && value != userId {
			userLocation, err := r.IRedisRepository.GetUserLocation(value)
			if err != nil {
				return nil, err
			}
			isNear, err := checkCollaborateArea(myLocation, userLocation, circle)
			if err != nil {
				return nil, err
			}
			if isNear {
				result.UserId = append(result.UserId, value)
			}
		}
	}
	return result, nil
}

func (r *room) IsMemberReady(roomId string) (bool, error) {
	roomMemberLimitValue, err := r.IRedisRepository.GetRoomMemberLimit(roomId)
	if err != nil {
		return false, err
	}
	roomMemberValue, err := r.IRedisRepository.GetRoomMember(roomId)
	if err != nil {
		return false, err
	}
	if roomMemberLimitValue.Value != string(len(roomMemberValue.UserId)) {
		return false, fmt.Errorf("invailed Value")
	}
	return true, nil
}

const (
	EQUATORIAL_RADIUS    = 6378137.0            // 赤道半径 GRS80
	POLAR_RADIUS         = 6356752.314          // 極半径 GRS80
	ECCENTRICITY         = 0.081819191042815790 // 第一離心率 GRS80
)

func checkCollaborateArea(myLocation, userLocation *db.UserLocation, circle float64) (bool, error) {
	dx := degree2radian(myLocation.Longitude - userLocation.Longitude)
	dy := degree2radian(myLocation.Latitude - userLocation.Latitude)
	my := degree2radian((myLocation.Latitude + userLocation.Latitude) / 2)

	W := math.Sqrt(1 - (Power2(ECCENTRICITY) * Power2(math.Sin(my)))) // 卯酉線曲率半径の分母
	m_numer := EQUATORIAL_RADIUS * (1 - Power2(ECCENTRICITY))         // 子午線曲率半径の分子

	M := m_numer / math.Pow(W, 3) // 子午線曲率半径
	N := EQUATORIAL_RADIUS / W    // 卯酉線曲率半径

	d := math.Sqrt(Power2(dy*M) + Power2(dx*N*math.Cos(my)))
	if d <= circle {
		return true, nil
	}
	return false, nil
}

func degree2radian(x float64) float64 {
	return x * math.Pi / 180
}

func Power2(x float64) float64 {
	return math.Pow(x, 2)
}