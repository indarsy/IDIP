package middleware

import (
	"errors"
	"strings"
	"time"
	"videodb/be/config"
	"videodb/be/utils"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt"
)

type Claims struct {
	UserID   uint   `json:"userId"`
	Username string `json:"username"`
	Role     string `json:"role"`
	jwt.StandardClaims
}

// JWT认证中间件
func JWTAuth() gin.HandlerFunc {
	return func(c *gin.Context) {
		token := c.GetHeader("Authorization")
		if token == "" {
			utils.Error(c, utils.ErrUnauthorized)
			c.Abort()
			return
		}

		// 移除Bearer前缀
		token = strings.Replace(token, "Bearer ", "", 1)

		// 解析token
		claims, err := ParseToken(token)
		if err != nil {
			utils.Error(c, utils.ErrUnauthorized)
			c.Abort()
			return
		}

		// 将用户信息存储到上下文
		c.Set("userId", claims.UserID)
		c.Set("username", claims.Username)
		c.Set("role", claims.Role)

		c.Next()
	}
}

// 生成token
func GenerateToken(userID uint, username, role string) (string, error) {
	claims := Claims{
		UserID:   userID,
		Username: username,
		Role:     role,
		StandardClaims: jwt.StandardClaims{
			ExpiresAt: time.Now().Add(config.GlobalConfig.JWT.ExpireTime).Unix(),
			IssuedAt:  time.Now().Unix(),
			Issuer:    "video-system",
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString([]byte(config.GlobalConfig.JWT.Secret))
}

// 解析token
func ParseToken(tokenString string) (*Claims, error) {
	token, err := jwt.ParseWithClaims(tokenString, &Claims{}, func(token *jwt.Token) (interface{}, error) {
		return []byte(config.GlobalConfig.JWT.Secret), nil
	})

	if err != nil {
		return nil, err
	}

	if claims, ok := token.Claims.(*Claims); ok && token.Valid {
		return claims, nil
	}

	return nil, errors.New("invalid token")
}
