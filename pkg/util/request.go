package util

import "github.com/gin-gonic/gin"

func BindAndValidate[T any](c *gin.Context) (T, error) {
	var req T
	if err := c.ShouldBind(&req); err != nil {
		return req, err
	}
	return req, nil
}
