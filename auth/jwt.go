package auth

import (
	"gitee.com/saltlamp/im-service/config"
	"gitee.com/saltlamp/im-service/syserr"
	"github.com/dgrijalva/jwt-go"
	"time"
)

const (
	AUTH_HEADER = "AUTH_TOKEN"
)

type JWT struct {
	AtNum    string
	Type     string
	Duration *time.Duration
}

type TokenType string

const (
	// admin
	TokenTypeAdmin string = "admin"
	// visitor
	TokenTypeVisitor string = "visitor"
	// user
	TokenTypeUser string = "user"
)

// create token
func CreateToken(JWT *JWT) (string, error) {
	if JWT.Duration == nil {
		return "", syserr.NewTokenAuthError("Duration is null")
	}
	if JWT.Type == "" {
		// default
		JWT.Type = TokenTypeUser
	}
	token := jwt.New(jwt.SigningMethodHS256)
	claims := make(jwt.MapClaims)
	claims["atNum"] = JWT.AtNum
	claims["type"] = JWT.Type
	claims["exp"] = time.Now().Add(*JWT.Duration).Unix()
	claims["iat"] = time.Now().Unix()
	token.Claims = claims
	if tokenString, err := token.SignedString([]byte(config.SignKey)); err == nil {
		return tokenString, nil
	} else {
		return "", syserr.NewTokenAuthError(err.Error())
	}
}

//  get id from by token
func GetJWT(token string) (*JWT, error) {
	tokenInfo, err := ParseToken(token)
	if err == nil {
		if claims, ok := tokenInfo.Claims.(jwt.MapClaims); ok && tokenInfo.Valid {
			return &JWT{
				AtNum: claims["atNum"].(string),
				Type:  claims["type"].(string),
			}, nil
		}
	}
	return nil, syserr.NewTokenAuthError(err.Error())
}

// check token
func ParseToken(s string) (*jwt.Token, error) {
	fn := func(token *jwt.Token) (interface{}, error) {
		return []byte(config.SignKey), nil
	}
	token, err := jwt.Parse(s, fn)
	if err != nil {
		err = syserr.NewTokenAuthError(err.Error())
	}
	return token, err
}
