package middleware

import (
	"demoApp/server/utils"
	"demoApp/server/utils/cache"
	"demoApp/server/utils/jwt_auth"
	"fmt"
	"github.com/dgrijalva/jwt-go"
	"github.com/dgrijalva/jwt-go/request"
	"github.com/julienschmidt/httprouter"
	"goframework/config"
	"goframework/gLog"
	"net/http"
)

func AuthorizationVerify(handle httprouter.Handle) httprouter.Handle {

	auth := func(w http.ResponseWriter, r *http.Request, param httprouter.Params) {

		token, err := request.ParseFromRequest(r, request.OAuth2Extractor, func(token *jwt.Token) (i interface{}, e error) {
			return []byte(config.JWTAuthSection.Key("salt").MustString(utils.SALT)), nil
		}, request.WithClaims(&jwt_auth.UserClaim{}))

		if err != nil || !token.Valid {
			goto last
		}

		if claim, ok := token.Claims.(*jwt_auth.UserClaim); ok {
			// 缓存判断该token
			b, err := cache.GetKeyFromCache(claim.Uuid)

			if err != nil || string(b) != token.Raw {
				gLog.LOG_INFO(fmt.Sprintf("can't found token  %s in cache ", string(b)))
				goto last
			}
			// 判断角色
			fmt.Println(claim.Role, r.URL.String())
			if config.PolicyEnforce.Enforce(claim.Role, r.URL.String(), r.Method) == false {
				gLog.LOG_INFO(fmt.Sprintf("user %s  with role %s "+
					"cann't access url %s", claim.Uuid, claim.Role, r.URL.String()))
				//goto last
				http.Error(w, http.StatusText(http.StatusForbidden), http.StatusForbidden)
				return
			}

			r.Header.Set("USER_UUID", claim.Uuid)
			r.Header.Set("USER_ROLE", claim.Role)
			handle(w, r, param)
			return
		}
	last:
		// 前端 如何判断 401
		http.Error(w, http.StatusText(http.StatusUnauthorized), http.StatusUnauthorized)

	}

	return httprouter.Handle(auth)
}
