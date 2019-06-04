package sender

import (
	"eyeSky/modules/transport/model"
	log "github.com/sirupsen/logrus"
	"github.com/toolkits/container/list"
)

const (
	DefaultSendQueueMaxSize = 10240000 //10.24w
)

var (
	s        sender
	nsqQueue = list.NewSafeListLimited(DefaultSendQueueMaxSize)
)

// send data to nsq
type sender struct{}

func Sender() sender {
	return s
}


func PushToNsqSendQueue(items []model.MetricValue) {
	log.Debugf("PushToNsqSendQueue: %d items", len(items))
	s.PushToNsqSendQueue(items)
}


func (s *sender) PushToNsqSendQueue(items []model.MetricValue) {

	for _, item := range items {

		isSuccess := nsqQueue.PushFront(item)
		if !isSuccess {
			log.Warn("nsqQueue.PushFront fail, please check")
		}
	}
}
