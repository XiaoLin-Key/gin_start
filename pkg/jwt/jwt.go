package jwt

import (
	"errors"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/spf13/viper"
)

var mySecret = []byte("Lin")

type MyClaims struct {
	UserID   int64  `json:"user_id"`
	Username string `json:"username"`
	jwt.RegisteredClaims
}

// GenToken 生成 JWT
func GenToken(userID int64, username string) (string, error) {
	claims := MyClaims{
		userID,
		username,
		jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(time.Hour * time.Duration(viper.GetInt64("auth.jwt_expire")))), // 过期时间
			Issuer:    "gin-start",                                                                                      // 签发人
		},
	}
	// 使用指定的签名方法创建签名对象
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	// 使用指定的secret签名并获得完整的编码后的字符串token
	return token.SignedString(mySecret)
}

// ParseToken 解析 JWT
func ParseToken(tokenStr string) (*MyClaims, error) {
	var mc = new(MyClaims)
	token, err := jwt.ParseWithClaims(tokenStr, mc, func(token *jwt.Token) (interface{}, error) {
		return mySecret, nil
	})
	if err != nil {
		return nil, err
	}
	if token.Valid {
		return mc, nil
	}
	return nil, errors.New("invalid token")
}

// 双token 模式，使用时记得在ParseToken中加上过期可以继续使用的逻辑
/*// GenToken 生成 Access Token 和 Refresh Token
func GenToken(userID int64, username string) (aToken, rToken string, err error) {
	// 1. Access Token: 短效，用于身份验证
	atClaims := MyClaims{
		userID,
		username,
		jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(time.Hour * 2)), // 2小时过期
			Issuer:    "bluebell",
		},
	}
	aToken, err = jwt.NewWithClaims(jwt.SigningMethodHS256, atClaims).SignedString(mySecret)

	// 2. Refresh Token: 长效，仅用于刷新 Access Token，通常不包含业务数据
	rtClaims := jwt.RegisteredClaims{
		ExpiresAt: jwt.NewNumericDate(time.Now().Add(time.Hour * 24 * 7)), // 7天过期
		Issuer:    "bluebell",
	}
	rToken, err = jwt.NewWithClaims(jwt.SigningMethodHS256, rtClaims).SignedString(mySecret)

	return
}

// RefreshToken 刷新 AccessToken
// 逻辑：验证 rToken，如果没问题，重新生成一对新的 Token
func RefreshToken(aToken, rToken string) (newAToken, newRToken string, err error) {
	// 1. 验证 Refresh Token 是否有效
	_, err = jwt.Parse(rToken, func(token *jwt.Token) (interface{}, error) {
		return mySecret, nil
	})
	if err != nil {
		return "", "", err
	}

	// 2. 从旧的 aToken 中解析出用户信息（注意：即便过期了也能解析出数据）
	var mc = new(MyClaims)
	_, err = jwt.ParseWithClaims(aToken, mc, func(token *jwt.Token) (interface{}, error) {
		return mySecret, nil
	})

	// 💡 这里的逻辑是：如果 aToken 仅仅是过期了，我们依然拿它的 UserID 重新签发
	v, _ := err.(*jwt.ValidationError)
	if v.Errors == jwt.ValidationErrorExpired {
		return GenToken(mc.UserID, mc.Username)
	}

	return "", "", errors.New("aToken 校验失败")
}
*/
