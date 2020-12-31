package impls

import (
	"context"
	_ "github.com/alexbrainman/odbc"
	_ "github.com/denisenkom/go-mssqldb"
	_ "github.com/go-sql-driver/mysql"
	_ "github.com/lib/pq"
	//_ "github.com/mattn/go-oci8"
	//_"github.com/go-goracle/goracle"
	"time"

	//_ "gopkg.in/goracle.v2"
	"database/sql"
	"github.com/astaxie/beedb"
	"github.com/astaxie/beego"
	"github.com/astaxie/beego/logs"
	"github.com/pkg/errors"
	"strconv"
)

type RDBMS struct {
	enabled    bool
	concurrent int
	maxquesize int
	dbtype     string
	connStr    string
	msgChan    chan *StoreMsg
	clients    []*sql.DB
	running    bool
}

var RDBMSSrv RDBMS

func (this *RDBMS) Status() map[string]interface{} {
	status := make(map[string]interface{})

	status["chan_len"] = len(this.msgChan)
	status["client_len"] = len(this.clients)
	status["concurrent"] = this.concurrent
	status["conn_string"] = this.connStr
	status["max_que_size"] = this.maxquesize

	status["running"] = this.running
	status["enabled"] = this.enabled
	status["dbtype"] = this.dbtype

	return status
}

func (this *RDBMS) initialize() error {
	var err error

	this.enabled, err = beego.AppConfig.Bool("rdbms::enabled")
	if err != nil {
		return errors.New("get rdbms::enabled parameter failed, " + err.Error())
	}

	if this.enabled {

		this.concurrent, err = beego.AppConfig.Int("rdbms::concurrent")
		if err != nil {
			logs.Error("get rdbms::concurrent parameter failed, ", err.Error())
			logs.Info("parameter rdbms::concurrent set default 10")
			this.concurrent = 10
		}

		if this.concurrent < 1 {
			logs.Info("parameter rdbms::concurrent set default 10")
			this.concurrent = 10
		}

		this.maxquesize, err = beego.AppConfig.Int("rdbms::maxquesize")
		if err != nil {
			logs.Error("get rdbms::maxquesize parameter failed, ", err.Error())
			logs.Info("parameter rdbms::maxquesize set default 100")
			this.maxquesize = 100
		}

		if this.maxquesize < 1 {
			logs.Info("parameter rdbms::maxquesize set default 1000")
			this.maxquesize = 1000
		}

		this.dbtype = beego.AppConfig.String("rdbms::dbtype")

		this.connStr = beego.AppConfig.String("rdbms::connStr")

	}

	this.msgChan = make(chan *StoreMsg, this.maxquesize+10)

	return nil
}

func (this *RDBMS) Start() error {
	err := this.initialize()
	if err != nil {
		return err
	}

	if !this.enabled {
		return nil
	}

	ok, _ := beego.AppConfig.Bool("rdbms::debug")
	if ok {
		beedb.OnDebug = true
	}

	this.clients = make([]*sql.DB, this.concurrent)

	db, err := sql.Open(this.dbtype, this.connStr)
	if err != nil {
		return errors.New("Start RDBMS store service failed, " + err.Error())
	}

	defer db.Close()

	db.SetMaxIdleConns(1)
	db.SetMaxOpenConns(1)

	//logs.Trace("db conn timeout:",timeout)

	connCtx, connCancel := context.WithTimeout(context.Background(), time.Duration(5)*time.Second)
	defer connCancel()

	err = db.PingContext(connCtx)

	if err != nil {
		return errors.New("Start RDBMS store service failed, " + err.Error())
	}

	for i := 0; i < this.concurrent; i++ {
		go this.run(i)
	}

	logs.Info("RDBMS store service started.")

	this.running = true
	return nil

}

func (this *RDBMS) PutMsg(msg *StoreMsg) error {
	if !this.running {
		return errors.New("RDBMS store service has not running")
	}

	if len(this.msgChan) >= this.maxquesize {
		errMsg := "RDBMS store service max queue size exceed, size: " + strconv.Itoa(this.maxquesize)
		logs.Error(errMsg)
		return errors.New(errMsg)
	}

	this.msgChan <- msg
	return nil
}

func (this *RDBMS) getConn(pos int) (*beedb.Model, error) {
	db, err := sql.Open(this.dbtype, this.connStr)
	if err != nil {
		return nil, errors.New("RDBMS store service connection create failed, " + err.Error())
	}

	db.SetMaxIdleConns(1)
	db.SetMaxOpenConns(1)

	db.SetConnMaxLifetime(10 * time.Minute)

	connCtx, connCancel := context.WithTimeout(context.Background(), time.Duration(5)*time.Second)
	defer connCancel()

	err = db.PingContext(connCtx)

	if err != nil {
		return nil, errors.New("RDBMS store service connect failed, " + err.Error())
	}

	this.clients[pos] = db

	var tdb beedb.Model
	switch this.dbtype {
	case "postgres", "pg":
		tdb = beedb.New(db, "pg")
	case "mssql", "sqlserver", "sybase":
		tdb = beedb.New(db, "mssql")
	default:
		tdb = beedb.New(db)
	}

	return &tdb, nil
}

func (this *RDBMS) run(pos int) {
	var err error
	var client *beedb.Model
	newConn := 0
	for {
		newConn = 0
		select {
		case msg, ok := <-this.msgChan:
			if !ok {
				return
			}

		RECONNECT:
			if client == nil {
				client, err = this.getConn(pos)
				if err != nil {
					logs.Error(err.Error())
					continue
				}
				newConn = 1
			}

			_, err = client.SetTable(msg.Series).InsertBatch(msg.Data)
			if err != nil {
				if newConn != 1 {
					logs.Warning("RDBMS store service insert rows failed, ", err.Error(), ", reconnect...")
					if this.clients[pos] != nil {
						this.clients[pos].Close()
					}
					this.clients[pos] = nil
					client = nil
					ReConnSleep()
					goto RECONNECT

				} else {
					errMsg := "RDBMS store service insert rows failed, " + err.Error()
					logs.Error(errMsg)
				}

			}
		}
	}

}
