package utils

import (
	"context"
	"encoding/json"
	"errors"
)

type UserIdInfo struct {
	UserId uint
	OpenId string
}

func GetCtxUserInfo(ctx context.Context) (*UserIdInfo, error) {
	openId, exist := ctx.Value("openId").(string)
	if !exist {
		return nil, errors.New("jwt openId not exist")
	}
	userId, exist := ctx.Value("userId").(json.Number)
	if !exist {
		return nil, errors.New("jwt userId not exist")
	}
	userIdInt, err := userId.Int64()
	if err != nil {
		return nil, errors.New("jwt userId not int64")
	}

	return &UserIdInfo{
		UserId: uint(userIdInt),
		OpenId: openId,
	}, nil
}
