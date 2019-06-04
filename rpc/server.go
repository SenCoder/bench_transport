/*
 * Copyright (c) 2018 PingAn. All rights reserved.
 */

package rpc

import (
	"eyeSky/modules/transport/conf"
	log "github.com/sirupsen/logrus"
	"net"
	"net/rpc"
)

func Start() {
	if !conf.Config().Rpc.Enabled {
		return
	}

	addr := conf.Config().Rpc.Listen

	tcpAddr, err := net.ResolveTCPAddr("tcp", addr)
	if err != nil {
		log.Fatalf("net.ResolveTCPAddr fail: %s", err)
	}

	listener, err := net.ListenTCP("tcp", tcpAddr)
	if err != nil {
		log.Fatalf("listen %s fail: %s", addr, err)
	} else {
		log.Info("rpc server start on ", addr)
	}

	rpc.Register(new(Transport))

	for {
		conn, err := listener.Accept()
		if err != nil {
			log.Warn("listener.Accept occur error:", err)
			continue
		}
		go rpc.ServeConn(conn)
	}
}
