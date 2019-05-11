package guid

import "github.com/rs/xid"

type guidUtil struct {}

type IGuidUtil interface {
	CreateGuid() string
}

func NewGuidUtil () IGuidUtil {
	return &guidUtil{}
}

func (util *guidUtil) CreateGuid() string {
	guid := xid.New()
	return guid.String()
}
