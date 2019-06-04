/*
 * Copyright (c) 2018 PingAn. All rights reserved.
 */

package http

import (

	"eyeSky/modules/transport/conf"
	"eyeSky/modules/transport/http/router"
	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
)

func Start() {

	httpCfg := conf.Config().Http

	if !httpCfg.Enabled {
		return
	}

	g := gin.New()
	router.Init(g)
	go g.Run(httpCfg.Listen)

	log.Infof("http server start on %s", httpCfg.Listen)
}
