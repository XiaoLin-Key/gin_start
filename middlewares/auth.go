package middlewares

import (
	"gin_start/pkg/jwt"
	"gin_start/result"
	"strings"

	"github.com/gin-gonic/gin"
)

const CtxUserID = "userID"

func JWTAuthMiddleware() func(c *gin.Context) {
	return func(c *gin.Context) {
		//	判断请求头是否有有效的token
		auth := c.Request.Header.Get("Authorization")
		if auth == "" {
			result.ResponseError(c, result.CodeNeedLogin)
			c.Abort()
			return
		}
		//按照空格分割token
		partd := strings.Split(auth, " ")
		if len(partd) != 2 || partd[0] != "Bearer" {
			result.ResponseError(c, result.CodeInvalidAuth)
			c.Abort()
			return
		}
		token := partd[1]

		//	解析token
		mc, err := jwt.ParseToken(token)
		if err != nil {
			/* 💡【双 Token 模式注释】:
			   如果前端发现 AccessToken 过期（比如后端返回一个特殊的 Code），
			   前端应该带着 RefreshToken 请求一个专门的刷新接口：/refresh_token

			   if errors.Is(err, jwt.ErrTokenExpired) {
			       result.ResponseErrorWithMsg(c, result.CodeTokenExpired, "token已过期，请刷新")
			       c.Abort()
			       return
			   }
			*/
			result.ResponseError(c, result.CodeInvalidAuth)
			c.Abort()
			return
		}
		//	将解析后的token信息存储到上下文
		c.Set(CtxUserID, mc.UserID)
		c.Next()
	}
}
