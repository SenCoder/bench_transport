/*
 * Copyright (c) 2018 PingAn. All rights reserved.
 */

package rpc

import (
	"eyeSky/modules/transport/model"
	"eyeSky/modules/transport/conf"
	"eyeSky/modules/transport/sender"
	log "github.com/sirupsen/logrus"
)

type Transport struct{}

func (Transport) Ping(req model.NullRpcRequest, resp *model.SimpleRpcResponse) error {
	resp.Code = model.RPC_STATUS_OK
	log.Info("Transport Ping:", resp.String())
	return nil
}

func (Transport) Update(args *model.RpcArgs, reply *model.RpcReply) error {
	log.Debug("Transport Update")
	//log.Debug(args.Metrics)
	return UpdateMetricValues(args.Metrics, reply, "rpc")
}


func UpdateMetricValues(metrics []model.MetricValue, reply *model.RpcReply, from string) error {

		//if md, err := v.ToMetaData(); err != nil {
		//	reply.Invalid ++
		//	continue
		//} else {
		//	items = append(items, md)
		//}
    items := []model.MetricValue{}

    for _, v := range metrics {
        items = append(items, v)
    }
	log.Infof("%d invalid items", reply.Invalid)

	cfg := conf.Config()
	if cfg.NSQ.Enabled {
		sender.PushToNsqSendQueue(items)
	}

	return nil
}
