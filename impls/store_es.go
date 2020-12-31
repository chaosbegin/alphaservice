package impls

import (
	"context"
	"github.com/astaxie/beego"
	"github.com/astaxie/beego/logs"
	"github.com/chilts/sid"
	"github.com/olivere/elastic"
	"github.com/pkg/errors"
	"strconv"
	"strings"
)

type Es struct {
	enabled    bool
	apiUrl     string
	concurrent int
	maxquesize int
	username   string
	password   string
	msgChan    chan *StoreMsg
	clients    []*elastic.Client
	running    bool
}

var EsSrv Es

func (this *Es) initialize() error {
	var err error
	this.clients = make([]*elastic.Client, 0)

	this.enabled, err = beego.AppConfig.Bool("es::enabled")
	if err != nil {
		return errors.New("get es::enabled parameter failed, " + err.Error())
	}

	if this.enabled {
		this.apiUrl = beego.AppConfig.String("es::api_url")
		this.concurrent, err = beego.AppConfig.Int("es::concurrent")
		if err != nil {
			logs.Error("get es::concurrent parameter failed, ", err.Error())
			logs.Info("parameter es::concurrent set default 10")
			this.concurrent = 10
		}

		if this.concurrent < 1 {
			logs.Info("parameter es::concurrent set default 10")
			this.concurrent = 10
		}

		this.maxquesize, err = beego.AppConfig.Int("es::maxquesize")
		if err != nil {
			logs.Error("get es::maxquesize parameter failed, ", err.Error())
			logs.Info("parameter es::maxquesize set default 1000")
			this.maxquesize = 1000
		}

		if this.maxquesize < 1 {
			logs.Info("parameter es::maxquesize set default 1000")
			this.maxquesize = 1000
		}

		this.username = beego.AppConfig.String("es::username")
		this.password = beego.AppConfig.String("es::password")

	}

	this.msgChan = make(chan *StoreMsg, this.maxquesize+10)

	return nil
}

func (this *Es) Start() error {
	err := this.initialize()
	if err != nil {
		return err
	}

	if !this.enabled {
		return nil
	}

	urls := strings.Split(this.apiUrl, ";")

	for i := 0; i < this.concurrent; i++ {
		if len(this.username) > 0 {
			client, err := elastic.NewClient(elastic.SetSniff(false), elastic.SetURL(urls...))
			if err != nil {
				return err
			}
			this.clients = append(this.clients, client)
		} else {
			client, err := elastic.NewClient(elastic.SetSniff(false), elastic.SetURL(urls...), elastic.SetBasicAuth(this.username, this.password))
			if err != nil {
				return err
			}
			this.clients = append(this.clients, client)
		}
	}

	for _, c := range this.clients {
		go this.run(c)
	}

	logs.Info("ELK service started.")
	this.running = true
	return nil

}

func (this *Es) PutMsg(msg *StoreMsg) error {
	if !this.running {
		return errors.New("ELK store service has not running")
	}
	if len(this.msgChan) >= this.maxquesize {
		errMsg := "ELK service max queue size exceed, size: " + strconv.Itoa(this.maxquesize)
		logs.Error(errMsg)
		return errors.New(errMsg)
	}

	this.msgChan <- msg
	return nil
}

func (this *Es) run(client *elastic.Client) {
	var err error
	ctx := context.Background()
	bulkRequest := client.Bulk()
	for {
		select {
		case msg, ok := <-this.msgChan:
			if !ok {
				return
			}

			for _, row := range msg.Data {
				req := elastic.NewBulkIndexRequest().Index(msg.Series).Type(msg.Series).Id(sid.Id()).Doc(row)
				bulkRequest = bulkRequest.Add(req)
			}
			_, err = bulkRequest.Do(ctx)
			if err != nil {
				logs.Error("ELK client bulk request failed, ", err.Error())
			}

		}
	}

}
