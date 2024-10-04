package api

import "github.com/gin-gonic/gin"

func genFailRes(err error) gin.H {
	return gin.H{
		"ok":   false,
		"err":  err.Error(),
		"body": nil,
	}
}
func genOkRes(body interface{}) gin.H {
	return gin.H{
		"ok":   true,
		"err":  nil,
		"body": body,
	}
}
