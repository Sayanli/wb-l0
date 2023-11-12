package main

import (
	"fmt"
	"os"

	"github.com/nats-io/stan.go"
	"github.com/sirupsen/logrus"
	"github.com/xlab/closer"
)

const (
	clusterID = "test-cluster"
	channel   = "order-channel"
	clientID  = "order-publisher"
)

type publisher struct {
	clusterID string
	clientID  string
	channel   string
	stanConn  stan.Conn
}

func newPub() *publisher {
	return &publisher{
		clusterID: clusterID,
		clientID:  clientID,
		channel:   channel,
	}
}

func (p *publisher) initConnection() error {
	sc, err := stan.Connect(p.clusterID, p.clientID)
	if err != nil {
		logrus.Error(err)
		return err
	}
	p.stanConn = sc

	return nil
}

func (p *publisher) cmdPublishJson() {
	for {
		fmt.Printf("Введите путь до json файла: ")

		var path string
		fmt.Scanln(&path)

		file, err := os.ReadFile(path)
		if err != nil {
			logrus.Info("Ошибка чтения файла")
			continue
		}

		err = p.stanConn.Publish(channel, file)
		if err != nil {
			logrus.Info("Не удалось отправить сообщение в канал")
			continue
		}
		logrus.Info("Сообщение успешно отправлено")
	}
}

func (p *publisher) Close() {
	logrus.Info("Закрытие соединения с NATS")
	p.stanConn.Close()
}

func main() {
	defer closer.Close()

	pub := newPub()

	pub.initConnection()
	closer.Bind(pub.Close)

	pub.cmdPublishJson()
}
