/*
 * Copyright (c) 2018 PingAn. All rights reserved.
 */

package sender

import (
	"encoding/json"
	"eyeSky/modules/transport/model"
	"eyeSky/modules/transport/conf"
	"fmt"
	"github.com/nsqio/go-nsq"
	log "github.com/sirupsen/logrus"
	"math/rand"
	"time"
)

type NsqProducer struct {
	producer *nsq.Producer
	inited   chan *nsq.Producer
}

func (p *NsqProducer) init() error {
	p.inited = make(chan *nsq.Producer)
	cfg := conf.Config()
	log.Debug(cfg.NSQ)

	go func() {
		rand.Seed(time.Now().UnixNano())
		addrs := cfg.NSQ.Addrs
		for i := range rand.Perm(len(addrs)) {
			producer, err := nsq.NewProducer(addrs[i], conf.NsqConfig())
			if err != nil {
				log.Error(err)
				continue
			}
			err = producer.Ping()
			if err != nil {
				log.Error(err)
				continue
			}
			p.inited <- producer
			break
		}
	}()

	select {
	case producer, ok := <- p.inited:
		if ok {
			p.producer = producer
			close(p.inited)
		}
	case <-time.After(time.Second * 5):
		close(p.inited)
		return fmt.Errorf("NsqProducer init fail")
	}
	return nil
}

type errorMessage struct {
	topic string
	body  [][]byte
	err   error
}

// 处理发送失败的逻辑
func (p *NsqProducer) Publish(msgs []model.MetricValue) {
	body := make([][]byte, len(msgs))
	for i, v := range msgs {
		bytes, _ := json.Marshal(v)
		body[i] = bytes
	}

	cfg := conf.Config()

	//nsq.MultiPublish()

	for _, topic := range cfg.NSQ.Topics {
		err := p.producer.MultiPublish(topic, body)
		if err != nil {
			log.WithError(err).Debug("Publish error, go to errorHandle")
			p.errorHandle(errorMessage{topic, body, err})
		}
	}
}

func (p *NsqProducer) errorHandle(errMsg errorMessage) {
	cfg := conf.Config()
	for i := 0; i < cfg.NSQ.Retry; i++ {
		err := p.producer.MultiPublish(errMsg.topic, errMsg.body)
		if err == nil {
			return
		}
	}
	log.Error(errMsg.err)
	// todo: 添加失败统计
	return
}

// 启动则开启死循环从 Queue 中读取数据，读到数据则发送给 NSQ
func (p *NsqProducer) Fetch() {
	cfg := conf.Config()
	batch := cfg.NSQ.Batch
	log.Debugf("fetching data from queue, batch=%d", batch)

	for {
		items := nsqQueue.PopBackBy(batch)

		if len(items) > 0 {
			log.Debugf("fetched data: %d", len(items))
			messages := make([]model.MetricValue, 0)

			for _, v := range items {
				if item, ok := v.(model.MetricValue); ok {
					messages = append(messages, item)
				} else {
					continue
				}
			}
			p.Publish(messages)
		}
	}
}

func StartSendTask() {
	cfg := conf.Config()

	if cfg.NSQ.Enabled {
		p := NsqProducer{}
		if err := p.init(); err != nil {
			log.Fatal(err)
		}
		go p.Fetch()
	}

	// init semaphore
	//judgeConcurrent := cfg.NSQ.MaxConns

	//if tsdbConcurrent < 1 {
	//	tsdbConcurrent = 1
	//}
	//
	//if judgeConcurrent < 1 {
	//	judgeConcurrent = 1
	//}
	//
	//if graphConcurrent < 1 {
	//	graphConcurrent = 1
	//}

	// init send go-routines
	//for node := range cfg.Judge.Cluster {
	//	queue := JudgeQueues[node]
	//	go forward2JudgeTask(queue, node, judgeConcurrent)
	//}

	//if cfg.Tsdb.Enabled {
	//	go forward2TsdbTask(tsdbConcurrent)
	//}
}
