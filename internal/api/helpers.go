package api

import "github.com/gin-gonic/gin"

func genFailBody(input interface{}) gin.H {
	err := input
	if realErr, ok := err.(error); ok {
		err = realErr.Error()
	}
	return gin.H{
		"ok":   false,
		"err":  err,
		"body": nil,
	}
}
func genOkBody(body interface{}) gin.H {
	return gin.H{
		"ok":   true,
		"err":  nil,
		"body": body,
	}
}
