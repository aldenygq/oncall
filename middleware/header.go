package middleware

import (
	//"encoding/base64"
	"fmt"
	"net/http"
	//"oncall/middleware"

	//"net/url"
	//"strconv"
	"time"
	//"oncall/models"
	"github.com/gin-gonic/gin"
)

// NoCache is a middleware function that appends headers
// to prevent the client from caching the HTTP response.
func NoCache(c *gin.Context) {
	c.Header("Cache-Control", "no-cache, no-store, max-age=0, must-revalidate, value")
	c.Header("Expires", "Thu, 01 Jan 1970 00:00:00 GMT")
	c.Header("Last-Modified", time.Now().UTC().Format(http.TimeFormat))
	c.Next()
}

// Options is a middleware function that appends headers
// for options requests and aborts then exits the middleware
// chain and ends the request.
func Options(c *gin.Context) {
	if c.Request.Method != "OPTIONS" {
		c.Next()
	} else {
		c.Header("Access-Control-Allow-Origin", "*")
		c.Header("Access-Control-Allow-Methods", "GET,POST,PUT,PATCH,DELETE,OPTIONS")
		c.Header("Access-Control-Allow-Headers", "authorization, origin, content-type, accept")
		c.Header("Allow", "HEAD,GET,POST,PUT,PATCH,DELETE,OPTIONS")
		c.Header("Content-Type", "application/json")
		c.AbortWithStatus(200)
	}
}

// Secure is a middleware function that appends security
// and resource access headers.
func Secure(c *gin.Context) {
	c.Header("Access-Control-Allow-Origin", "*")
	c.Header("X-Content-Type-Options", "nosniff")
	c.Header("X-XSS-Protection", "1; mode=block")
	if c.Request.TLS != nil {
		c.Header("Strict-Transport-Security", "max-age=31536000")
	}
}
func NoRoute(c *gin.Context) {
	ctx := Context{Ctx: c}
	path := c.Request.URL.Path
	method := c.Request.Method
	
	ctx.Response(HTTP_NOT_FOUND_CODE, fmt.Sprintf("%s %s not found", method, path), nil)
}



//判断是否https
func IsHttps(c *gin.Context) bool {
	if c.GetHeader("X-Forwarded-Proto") =="https" || c.Request.TLS!=nil{
		return true
	}
	return false
}

