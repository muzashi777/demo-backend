package demo

import (
	db "demo/_config"
	"demo/controller"
	"fmt"

	"github.com/gin-gonic/gin"
)

func StartServer() {
	fmt.Println("Server starting tag...")
	r := gin.Default()

	resource, err := db.CreateResource()
	if err != nil {

		fmt.Println("err_db:", err)
	}
	v1 := r.Group("/v1")
	{
		v1.GET("/demo", controller.DemoV1(resource))
	}
	r.Run(":4488")
}
