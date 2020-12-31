package impls

import (
	"alphawolf.com/alpha/util"
	"context"
	"github.com/astaxie/beego"
	"github.com/astaxie/beego/logs"
	"github.com/pkg/errors"
	"github.com/segmentio/kafka-go"
	"strconv"
	"strings"
)

type Kafka struct {
	enabled    bool
	brokers    string
	topic      string
	concurrent int
	maxquesize int
	msgChan    chan *StoreMsg
	clients    []*kafka.Writer
	running    bool
}

var KafkaSrv Kafka

func (this *Kafka) initialize() error {
	var err error
	this.clients = make([]*kafka.Writer, 0)

	this.enabled, err = beego.AppConfig.Bool("kafka::enabled")
	if err != nil {
		return errors.New("get kafka::enabled parameter failed, " + err.Error())
	}

	if this.enabled {
		this.brokers = beego.AppConfig.String("kafka::brokers")
		this.topic = beego.AppConfig.String("kafka::topic")
		this.concurrent, err = beego.AppConfig.Int("kafka::concurrent")
		if err != nil {
			logs.Error("get kafka::concurrent parameter failed, ", err.Error())
			logs.Info("parameter kafka::concurrent set default 10")
			this.concurrent = 10
		}

		if this.concurrent < 1 {
			logs.Info("parameter kafka::concurrent set default 10")
			this.concurrent = 10
		}

		this.maxquesize, err = beego.AppConfig.Int("kafka::maxquesize")
		if err != nil {
			logs.Error("get kafka::maxquesize parameter failed, ", err.Error())
			logs.Info("parameter kafka::maxquesize set default 1000")
			this.maxquesize = 1000
		}

		if this.maxquesize < 1 {
			logs.Info("parameter kafka::maxquesize set default 1000")
			this.maxquesize = 1000
		}

	}

	this.msgChan = make(chan *StoreMsg, this.maxquesize+10)

	return nil
}

func (this *Kafka) Start() error {
	err := this.initialize()
	if err != nil {
		return err
	}

	if !this.enabled {
		return nil
	}

	brks := strings.Split(this.brokers, ";")

	for i := 0; i < this.concurrent; i++ {
		client := kafka.NewWriter(kafka.WriterConfig{
			Brokers:  brks,
			Topic:    this.topic,
			Balancer: &kafka.Hash{},
		})
		this.clients = append(this.clients, client)
	}

	for _, c := range this.clients {
		go this.run(c)
	}

	logs.Info("Kafka service started.")
	this.running = true

	return nil

}

func (this *Kafka) PutMsg(msg *StoreMsg) error {
	if !this.running {
		return errors.New("Kafka store service has not running")
	}
	if len(this.msgChan) >= this.maxquesize {
		errMsg := "Kafka service max queue size exceed, size: " + strconv.Itoa(this.maxquesize)
		logs.Error(errMsg)
		return errors.New(errMsg)
	}

	this.msgChan <- msg
	return nil
}

func (this *Kafka) run(client *kafka.Writer) {
	for {
		select {
		case msg, ok := <-this.msgChan:
			if !ok {
				return
			}

			data, err := util.JsonIter.Marshal(msg.Data)
			if err != nil {
				errMsg := "Marshal data to json failed, " + err.Error()
				logs.Error(errMsg)
				continue
			}

			err = client.WriteMessages(context.Background(),
				kafka.Message{
					Key:   []byte(msg.Series),
					Value: data,
				},
			)

			if err != nil {
				errMsg := "Kafka write data failed, " + err.Error()
				logs.Error(errMsg)
			}

		}
	}

}
