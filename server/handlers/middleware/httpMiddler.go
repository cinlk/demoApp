package middleware

// 判断token 是否有效
// 判断角色 是否能访问该api

//func AuthorizationVerify(handle http.Handler) http.Handler {
//
//	auth := func(w http.ResponseWriter, r *http.Request) {
//
//		token, err := request.ParseFromRequest(r, request.OAuth2Extractor, func(token *jwt.Token) (i interface{}, e error) {
//			return jwt_auth.DefaultUserClaim.Salt, nil
//		}, request.WithClaims(&jwt_auth.UserClaim{}))
//
//		if err != nil || !token.Valid {
//			goto last
//		}
//
//		if claim, ok := token.Claims.(*jwt_auth.UserClaim); ok {
//			// 缓存判断该token
//			b, err := cache.GetKeyFromCache(claim.Uuid)
//			if err != nil || string(b) != token.Raw {
//				goto last
//			}
//			// 判断角色
//			if config.PolicyEnforce.Enforce(claim.Role, r.URL.String(), r.Method) == false {
//				goto last
//			}
//
//			r.Header.Set("USER_UUID", claim.Uuid)
//			r.Header.Set("USER_ROLE", claim.Role)
//			handle.ServeHTTP(w, r)
//			return
//		}
//	last:
//		http.Error(w, http.StatusText(http.StatusUnauthorized), http.StatusUnauthorized)
//
//	}
//
//	return http.HandlerFunc(auth)
//}
