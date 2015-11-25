package service

import (
	"cubar.com/lib/authentication"
	"cubar.com/model"
	"encoding/json"
	jwt "github.com/dgrijalva/jwt-go"
	"net/http"
)

func Login(requestLogin *model.Login) (int, []byte) {

	authBackend := authentication.InitJWTAuthenticationBackend()

	if authBackend.Authenticate(requestLogin) {
		token, err := authBackend.GenerateToken(requestLogin.Account)
		if err != nil {
			return http.StatusInternalServerError, []byte("")
		} else {
			response, _ := json.Marshal(authentication.TokenAuthentication{token})
			return http.StatusOK, response
		}
	}

	return http.StatusUnauthorized, []byte("")
}

func RefreshToken(requestLogin *model.Login) []byte {

	authBackend := authentication.InitJWTAuthenticationBackend()

	token, err := authBackend.GenerateToken(requestLogin.Account)
	if err != nil {
		panic(err)
	}

	response, err := json.Marshal(authentication.TokenAuthentication{token})
	if err != nil {
		panic(err)
	}
	return response
}

func Logout(req *http.Request) error {

	authBackend := authentication.InitJWTAuthenticationBackend()
	tokenRequest, err := jwt.ParseFromRequest(req, func(token *jwt.Token) (interface{}, error) {
		return authBackend.PublicKey, nil
	})

	if err != nil {
		return err
	}

	tokenString := req.Header.Get("Authorization")

	return authBackend.Logout(tokenString, tokenRequest)
}
