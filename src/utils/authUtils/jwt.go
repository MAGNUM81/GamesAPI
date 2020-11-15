package authUtils

import "github.com/dgrijalva/jwt-go"

var superSecretString = []byte("golang.gang.gang")

type UserClaims struct {
	UserId uint64 `json:"user_id"`
	jwt.StandardClaims
}

//jwt utils inspired from https://tabvn.medium.com/authenticate-jwt-in-go-graphql-d71db976f71c

func JwtDecode(token string) (*jwt.Token, error) {
	return jwt.ParseWithClaims(token, &UserClaims{}, func(token *jwt.Token) (interface{}, error) {
		return superSecretString, nil
	})
}

func JwtCreate(userID uint64, expiredAt int64) string {
	claims := UserClaims{
		userID,
		jwt.StandardClaims{
			ExpiresAt: expiredAt,
			Issuer:    "GamesAPI",
		},
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	ss, _ := token.SignedString(superSecretString)
	return ss
}