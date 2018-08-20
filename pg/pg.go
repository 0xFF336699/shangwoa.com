package pg

import (
	"database/sql"

	"time"

	_ "github.com/lib/pq"
	"shangwoa.com/io2"
)

func GetClient(params string) (*sql.DB, error) {
	return sql.Open("postgres", params)
}

type Pg struct {
	//Pg        *sql.DB
	//Client    *sql.DB
	Connector *io2.Connector
}

func (this *Pg) GetClient(reconnect bool) *sql.DB {
	c, ch := this.Connector.GetClient(reconnect)
	if c == nil {
		c = <-ch
	}
	client := c.(*sql.DB)
	return client
}

/**

pg, err := common.CreatePg(Model.Def.CommConf.HBPg.String)
if err != nil {
	log2.Panic(err)
}
Model.Def.CommRemoteCollection.CnPg.Pg = pg
maps, err := gosqlmf.QueryAll(pg.GetClient(false), "select * from shortcode_media_vo limit 2")
fmt.Println(maps)
*/
func CreatePg(params string) (*Pg, error) {
	p := &Pg{
		Connector: CreateConnector(params),
	}
	return p, nil
}

func CreateConnector(params string) *io2.Connector {
	alarmed := false
	return &io2.Connector{
		Connect: func() (interface{}, error) {
			return GetClient(params)
		}, OnError: func(e error) bool {
			return false
		}, OnMaxSeqAlarm: func(i int) {
			if !alarmed {
				alarmed = true
				// alarm
			}
		}, MaxSeq: 2000,
		RetryDelay: time.Second * 2,
	}
}
