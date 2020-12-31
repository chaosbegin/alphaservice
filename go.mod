module alphawolf.com/alphaservice

go 1.13

require (
	alphawolf.com/alpha/crypto/rsa v0.0.0-00010101000000-000000000000 // indirect
	alphawolf.com/alpha/util v0.0.0-00010101000000-000000000000
	alphawolf.com/alphaservice/service v0.0.0-00010101000000-000000000000
	github.com/Luzifer/go-openssl v2.0.0+incompatible
	github.com/Luzifer/go-openssl/v3 v3.1.0
	github.com/OwnLocal/goes v1.0.0 // indirect
	github.com/alexbrainman/odbc v0.0.0-20190102080306-cf37ce290779
	github.com/astaxie/beedb v0.0.0-20141221130223-1732292dfde4
	github.com/astaxie/beego v1.12.3
	github.com/chilts/sid v0.0.0-20190607042430-660e94789ec9
	github.com/denisenkom/go-mssqldb v0.0.0-20191001013358-cfbb681360f0
	github.com/go-sql-driver/mysql v1.5.0
	github.com/gomarkdown/markdown v0.0.0-20201113031856-722100d81a8e
	github.com/gorilla/websocket v1.4.1
	github.com/influxdata/influxdb-client-go/v2 v2.2.0 // indirect
	github.com/influxdata/influxdb1-client v0.0.0-20200827194710-b269163b24ab
	github.com/lib/pq v1.2.0
	github.com/mailru/easyjson v0.7.6 // indirect
	github.com/microcosm-cc/bluemonday v1.0.2
	github.com/olivere/elastic v6.2.26+incompatible
	github.com/pkg/errors v0.9.1
	github.com/satori/go.uuid v1.2.1-0.20181028125025-b2ce2384e17b
	github.com/segmentio/kafka-go v0.3.4
	github.com/shiena/ansicolor v0.0.0-20151119151921-a422bbe96644 // indirect
	github.com/siddontang/go v0.0.0-20180604090527-bdc77568d726 // indirect
	github.com/siddontang/ledisdb v0.0.0-20181029004158-becf5f38d373 // indirect
	github.com/yuin/goldmark v1.2.1
	golang.org/dl v0.0.0-20200302224518-306f3096cb2f // indirect
	golang.org/x/crypto v0.0.0-20191119213627-4f8c1d86b1ba
)

replace alphawolf.com/alpha/util => ../alpha/util

replace alphawolf.com/alpha/crypto/rsa => ../alpha/crypto/rsa

replace alphawolf.com/alphaservice/service => ./service
