package global

import (
	"context"
	"fmt"
	"go-klikdokter/helper/message"

	"github.com/go-kit/kit/auth/jwt"
	jwtgo "github.com/golang-jwt/jwt"
	"github.com/spf13/viper"
)

type JWTInfo struct {
	ActorName string `json:"-"`
	ActorUID  string `json:"-"`
}

func SetJWTInfoFromContext(ctx context.Context) (*JWTInfo, message.Message) {
	jwtInfo := JWTInfo{}
	skipJWTValidation := viper.GetBool("security.jwt.skip-validation")

	if !skipJWTValidation {
		token, _, err := new(jwtgo.Parser).ParseUnverified(fmt.Sprint(ctx.Value(jwt.JWTContextKey)), jwtgo.MapClaims{})
		if err != nil {
			return &jwtInfo, message.ErrNoAuth
		}

		if claims, ok := token.Claims.(jwtgo.MapClaims); ok {
			jwtInfo.ActorName = fmt.Sprintf("%v", claims["name"])
			jwtInfo.ActorUID = fmt.Sprintf("%v", claims["sub"])
			return &jwtInfo, message.SuccessMsg
		} else {
			return &jwtInfo, message.ErrNoAuth
		}

	} else {
		return &jwtInfo, message.SuccessMsg
	}
}
