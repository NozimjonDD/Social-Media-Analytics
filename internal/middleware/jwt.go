package middleware

import (
	"github.com/gin-gonic/gin"
	"net/http"
	"strings"
	"tgstat/internal/entity"
)

type TokenVerifier interface {
	VerifyToken(tokenString string) (*entity.User, error)
}

type Middleware struct {
	token TokenVerifier
}

func NewMiddleware(token TokenVerifier) *Middleware {
	return &Middleware{
		token: token,
	}
}

func (m *Middleware) JWT() gin.HandlerFunc {
	return func(c *gin.Context) {
		header := c.GetHeader("Authorization")
		headerSplit := strings.Split(header, " ")
		if len(headerSplit) != 2 {
			c.Status(http.StatusUnauthorized)
			c.Abort()
			return
		}

		token := headerSplit[1]

		user, err := m.token.VerifyToken(token)
		if err != nil {
			c.Status(http.StatusUnauthorized)
			c.Abort()
			return
		}

		c.Set("user", user)
		c.Next()
	}
}
