package middleware

import (
	"github.com/gin-gonic/gin"
	"net/http"
)

var (
	maxBodyBytes int64 = 25 << 20
)

func BodySizeMiddleware(c *gin.Context) {
	var w http.ResponseWriter = c.Writer
	c.Request.Body = http.MaxBytesReader(w, c.Request.Body, maxBodyBytes)

	c.Next()
}
