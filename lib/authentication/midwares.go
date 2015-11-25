package authentication

import (
	"fmt"
	jwt "github.com/dgrijalva/jwt-go"
	// "log"
	"cubar.com/lib/util"
	"cubar.com/model"
	"net/http"
)

func RequireTokenAuthentication(rw http.ResponseWriter, req *http.Request) (bool, *model.User) {

	var loggedUser *model.User = nil
	var result bool = false
	authBackend := InitJWTAuthenticationBackend()

	token, err := jwt.ParseFromRequest(req, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodRSA); !ok {
			return nil, fmt.Errorf("Unexpected signing method: %v", token.Header["alg"])
		} else {
			return authBackend.PublicKey, nil
		}
	})

	if err == nil && token.Valid && !authBackend.IsInBlacklist(req.Header.Get("Authorization")) {

		// 验证通过了，设置当前登陆用户
		result = true

		account := util.ParseString(token.Claims["sub"])
		loggedUser = model.GetByAccount(account)
	} else {
		rw.WriteHeader(http.StatusUnauthorized)
	}

	return result, loggedUser
}
