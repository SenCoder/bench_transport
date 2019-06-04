/*
 * Copyright (c) 2018 PingAn. All rights reserved.
 */

package main

import (
	"eyeSky/modules/transport/conf"
	"eyeSky/modules/transport/rpc"
	"eyeSky/modules/transport/sender"
	"flag"
	"fmt"
	log "github.com/sirupsen/logrus"
	"os"
	"runtime"
)

// Transport 默认发送数据到一个 NSQ, 若发送失败，会尝试切换到另一个 NSQ 结点
// 如果 NSQ 失联，数据将会一直缓存在 list 里面，内存消耗会激增，后续我们考虑使用带过期策略的缓存

func main() {

	log.Infof("GOMAXPROCS: %d", runtime.GOMAXPROCS(runtime.NumCPU()))

	cfg := flag.String("c", "cfg.json", "configure file")
	version := flag.Bool("v", false, "show version")
	debug := flag.Int("d", 0, "run in debug mode")
	flag.Parse()

	if *version {
		fmt.Println(conf.VERSION)
		os.Exit(0)
	}

	conf.LoadConfig(*cfg)
	if *debug == 1 {
		log.SetLevel(log.DebugLevel)
	} else {
		log.SetLevel(log.Level(conf.Config().LogLevel))
	}

	go rpc.Start()

	sender.StartSendTask()

	select {}

}
