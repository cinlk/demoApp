package jwt_auth

import (
	"demoApp/server/utils"
	"github.com/dgrijalva/jwt-go"
	"github.com/pkg/errors"
	"time"
)

type jwtParser interface {
	Decode(token string) (jwt.Claims, error)
	Encode() (string, error)
	SetSalt(salt []byte, claim jwt.Claims)
}

type DefaultTokenClaim struct {
	jwt.StandardClaims
	Salt        []byte    `json:"-"`
	TokenParser jwtParser `json:"-"`
}

// implement claims interface
func (t *DefaultTokenClaim) Valid() error {
	return t.StandardClaims.Valid()
}

func (t *DefaultTokenClaim) SetParser(parser jwtParser, claim jwt.Claims) {
	t.TokenParser = parser
	if t.Salt == nil {
		t.Salt = []byte(utils.SALT)
	}
	if claim == nil {
		claim = t
	}

	t.TokenParser.SetSalt(t.Salt, claim)
}

var defaultTokenClaim = DefaultTokenClaim{
	StandardClaims: jwt.StandardClaims{
		ExpiresAt: time.Now().Add(time.Hour).Unix(),
		IssuedAt:  time.Now().Unix(),
	},
}

func (t *DefaultTokenClaim) ValidateToken(token string) (*DefaultTokenClaim, bool) {
	if t.TokenParser == nil {
		return nil, false
	}
	ob, err := t.TokenParser.Decode(token)
	if err != nil {
		return nil, false
	}
	if claim, ok := ob.(*DefaultTokenClaim); ok {
		return claim, true
	}
	return nil, false
}

func SetSalt(salt []byte) func(claim *DefaultTokenClaim) {
	return func(claim *DefaultTokenClaim) {
		claim.Salt = salt
	}
}

func SetExpireTime(t time.Duration) func(claim *DefaultTokenClaim) {
	return func(claim *DefaultTokenClaim) {
		claim.ExpiresAt = time.Now().Add(time.Second * t).Unix()
	}
}

type DefaultTokenParser struct {
	Claim jwt.Claims
	Salt  []byte
}

func (t *DefaultTokenParser) Decode(token string) (jwt.Claims, error) {

	if t.Claim == nil {
		return nil, errors.New("not found claim")
	}

	ob, err := jwt.ParseWithClaims(token, t.Claim, func(token *jwt.Token) (i interface{}, e error) {
		return t.Salt, nil
	})
	if err != nil {
		return nil, err
	}
	if ob.Valid == false {
		return nil, errors.New("token is invalide")
	}

	return ob.Claims, nil
}

func (t *DefaultTokenParser) Encode() (string, error) {

	token := jwt.NewWithClaims(jwt.SigningMethodHS512, t.Claim)

	return token.SignedString(t.Salt)

}

func (t *DefaultTokenParser) SetSalt(salt []byte, claim jwt.Claims) {
	t.Salt = salt
	t.Claim = claim
}
