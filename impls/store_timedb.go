package impls

import (
	"errors"
	"fmt"
	"github.com/astaxie/beego"
	"github.com/astaxie/beego/logs"
	//"github.com/influxdata/influxdb-client-go/v2"
	_ "github.com/influxdata/influxdb1-client" // this is important because of the bug in go mod
	client "github.com/influxdata/influxdb1-client/v2"
	"strconv"
	"time"
)

type Timedb struct {
	enabled    bool
	concurrent int
	clientType string
	address    string
	database   string
	username   string
	password   string
	maxquesize int

	msgChan  chan *StoreMsg
	bpConfig client.BatchPointsConfig
	clients  []client.Client
	running  bool
}

var TimedbSrv Timedb

func (this *Timedb) initialize() error {
	var err error
	this.clients = make([]client.Client, 0)

	this.enabled, err = beego.AppConfig.Bool("timedb::enabled")
	if err != nil {
		return errors.New("get timedb::enabled parameter failed, " + err.Error())
	}

	this.concurrent, err = beego.AppConfig.Int("timedb::concurrent")
	if err != nil {
		return errors.New("Invalid timedb::concurrent parameter, " + err.Error())
	}

	if this.concurrent < 1 {
		this.concurrent = 1
	}

	this.maxquesize, err = beego.AppConfig.Int("timedb::maxquesize")
	if err != nil {
		logs.Error("Invalid parameter timedb::maxquesize")
		logs.Info("Parameter timedb::maxquesize set to default 1000")
		this.maxquesize = 1000
	}

	if this.maxquesize < 1 {
		logs.Info("parameter mysql::maxquesize set default 1000")
		this.maxquesize = 1000
	}

	this.msgChan = make(chan *StoreMsg, this.maxquesize+10)

	this.clientType = beego.AppConfig.String("timedb::client_type")

	this.database = beego.AppConfig.String("timedb::database")
	if len(this.database) < 1 {
		return errors.New("Invalid parameter timedb::database")
	}

	this.address = beego.AppConfig.String("timedb::address")
	if len(this.database) < 1 {
		return errors.New("Invalid parameter timedb::address")
	}

	this.username = beego.AppConfig.String("timedb::username")
	this.password = beego.AppConfig.String("timedb::password")

	return nil
}

func (this *Timedb) Start() (err error) {
	err = this.initialize()
	if err != nil {
		return err
	}

	for i := 0; i < this.concurrent; i++ {
		switch this.clientType {
		case "udp":
			client, err := client.NewUDPClient(client.UDPConfig{
				this.address,
				65535})

			if err != nil {
				return err
			}

			this.clients = append(this.clients, client)
		case "http":
			client, err := client.NewHTTPClient(client.HTTPConfig{
				Addr:     this.address,
				Username: this.username,
				Password: this.password,
			})

			if err != nil {
				return err
			}

			this.clients = append(this.clients, client)
		default:
			return errors.New("Unknown timedb client type " + this.clientType)

		}
	}

	for _, c := range this.clients {
		go this.run(c)
	}

	logs.Info("Timedb service started")

	this.running = true
	return nil
}

func (this *Timedb) run(client client.Client) {
	for {
		select {
		case msg, ok := <-this.msgChan:
			{
				if !ok {
					return
				}

				bp, err := this.getBps(msg)
				if err != nil {
					logs.Error("Timedb client get point failed, ", err.Error())
				}

				// Write the batch
				err = client.Write(*bp)
				if err != nil {
					logs.Error("Timedb client write point failed, ", err.Error())
				}

			}
		}
	}

}

func (this *Timedb) PutMsg(msg *StoreMsg) error {
	if !this.running {
		return errors.New("TimeDB store service has not running")
	}
	if len(this.msgChan) >= this.maxquesize {
		errMsg := "TimeDB service max queue size exceed, size: " + strconv.Itoa(this.maxquesize)
		logs.Error(errMsg)
		return errors.New(errMsg)
	}

	this.msgChan <- msg
	return nil
}

func (this *Timedb) getBps(msg *StoreMsg) (*client.BatchPoints, error) {
	bp, err := client.NewBatchPoints(client.BatchPointsConfig{
		Database:  this.database,
		Precision: "us",
	})

	if err != nil {
		return nil, errors.New("Create batchpoint failed, " + err.Error())
	}

	for _, row := range msg.Data {
		tags := make(map[string]string)
		if msg.Tags != nil && len(msg.Tags) > 0 {
			for _, t := range msg.Tags {
				val, ok := row[t]
				if ok && len(fmt.Sprint(val)) > 0 {
					tags[t] = fmt.Sprint(val)
				}
			}

		}

		if len(msg.Host) > 0 {
			tags["_host"] = msg.Host
		}

		pt, err := client.NewPoint(msg.Series, tags, row, time.Now().UTC())
		if err != nil {
			return nil, errors.New("Create point failed, " + err.Error())
		}
		bp.AddPoint(pt)
	}

	return &bp, nil
}
