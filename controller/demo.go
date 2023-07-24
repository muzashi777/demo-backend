package controller

import (
	db "demo/_config"

	"github.com/gin-gonic/gin"
)

//resource *db.Resource
func DemoV1(resource *db.Resource) func(c *gin.Context) {
	type Body struct {
		Username string `json:"username" bson:"username"`
	}
	return func(c *gin.Context) {

		c.JSON(200, gin.H{
			"status": map[string]interface{}{
				"code":    200,
				"message": "success",
			},
			"message": map[string]interface{}{
				"code":    200,
				"message": "success",
			},
			"data": "i am here",
		})
		return
	}
}
