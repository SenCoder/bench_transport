/*
 * Copyright (c) 2018 PingAn. All rights reserved.
 */
 
package router

import (
	"eyeSky/modules/api/router/middleware"
	"github.com/gin-gonic/gin"
)

func Init(router *gin.Engine) {
	router.Use(middleware.CORS())

	router.POST("/api/push", apiPushHandler)
	router.GET("/api/ping", apiPingHandler)
}
