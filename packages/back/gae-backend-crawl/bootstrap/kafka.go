package bootstrap

import (
	"gae-backend-crawl/constant/mq"
	kq "gae-backend-crawl/infrastructure/kafka"
	"github.com/zeromicro/go-zero/core/service"
)

func initKafkaConf(*Env) map[int]kq.KqConf {

	confMap := new(map[int]kq.KqConf)
	m := *confMap
	m = make(map[int]kq.KqConf)

	// 为 UnCleansingRepo 配置
	UnCleansingRepoGroup := kq.KqConf{
		ServiceConf: service.ServiceConf{
			Name: "gaeUnCleansingRepoService",
		},
		Brokers: []string{mq.KafkaDefaultLocalBroker},
		Group:   mq.UnCleansingRepoGroup,
		Topic:   mq.UnCleansingRepoTopic,
		Offset:  mq.FirstOffset,
		Conns:   1,
	}

	m[mq.UnCleansingRepoId] = UnCleansingRepoGroup

	// 为 UnCleansingUser 配置
	UnCleansingUserGroup := kq.KqConf{
		ServiceConf: service.ServiceConf{
			Name: "gaeUnCleansingUserService",
		},
		Brokers: []string{mq.KafkaDefaultLocalBroker},
		Group:   mq.UnCleansingUserGroup,
		Topic:   mq.UnCleansingUserTopic,
		Offset:  mq.FirstOffset,
		Conns:   1,
	}

	m[mq.UnCleansingUserId] = UnCleansingUserGroup

	return m
}

type KafkaConf struct {
	Conf map[int]kq.KqConf
}

func NewKafkaConf(e *Env) *KafkaConf {
	conf := initKafkaConf(e)
	return &KafkaConf{Conf: conf}
}
