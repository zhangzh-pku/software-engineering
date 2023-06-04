package main

import (
	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"github.com/zhangzh-pku/software-engineering/pkg/api"
)

func main() {
	router := gin.Default()

	router.Use(cors.Default())
	api.SetupRoutes(router)

	router.Run(":8080")
}
