package jwt_auth

import (
	"demoApp/server/utils/cache"
	"goframework/config"
	"goframework/gLog"
	"time"
)

var DefaultUserClaim *UserClaim

type CreateClaim interface {
	CreateDefaultToken(uuid, role string) (string, error)
}

type UserClaim struct {
	DefaultTokenClaim
	Uuid string `json:"uuid"`
	Role string `json:"role"`
}

func (u *UserClaim) CreateToken(uuid, role string, parser jwtParser, fs ...func(user *DefaultTokenClaim)) (string, error) {

	for _, f := range fs {
		f(&u.DefaultTokenClaim)
	}
	u.Uuid = uuid
	u.Role = role
	if parser != nil {
		u.SetParser(parser, u)
	} else {
		parser := &DefaultTokenParser{}
		u.SetParser(parser, u)
	}
	gLog.LOG_INFOF("the claim is %+v", u)
	return u.TokenParser.Encode()

}

func (u *UserClaim) CreateDefaultToken(uuid, role string) (string, error) {

	token, err := u.CreateToken(uuid, role, nil, SetSalt([]byte(config.JWTAuthSection.Key("salt").MustString(""))),
		SetExpireTime(time.Duration(config.JWTAuthSection.Key("expireTime").MustInt64(3600))))
	if err != nil {
		return "", err
	}
	// token 存入缓存
	if role != "anonymous" {
		err = cache.SetKeyInCache(uuid, []byte(token), 0)
		if err != nil {
			return "", err
		}
	}

	return token, nil

}

func init() {

	DefaultUserClaim = &UserClaim{
		defaultTokenClaim,
		"",
		"",
	}
}
