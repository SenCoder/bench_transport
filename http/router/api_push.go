/*
 * Copyright (c) 2018 PingAn. All rights reserved.
 */

package router

import (
	"encoding/json"
	"eyeSky/modules/base/model"
	"eyeSky/modules/transport/rpc"
	"github.com/gin-gonic/gin"
)

func apiPushHandler(c *gin.Context) {

	decoder := json.NewDecoder(c.Request.Body)
	var metrics []model.MetricValue
	err := decoder.Decode(&metrics)
	if err != nil {
		//http.Error(c.Writer, "decode error", http.StatusBadRequest)
		StdRender(c.Writer, nil, err)
		return
	}

	reply := &model.RpcReply{}
	rpc.UpdateMetricValues(metrics, reply, "http")
	RenderDataJson(c.Writer, reply)
}


func apiPingHandler(c *gin.Context) {

}
