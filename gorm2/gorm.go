package gorm2

import (
	"time"

	"github.com/jinzhu/gorm"
	_ "github.com/lib/pq"
	"shangwoa.com/io2"
)

type Gorm struct {
	Connector *io2.Connector
}

func (this *Gorm) GetClient(reconnect bool) *gorm.DB {
	c, ch := this.Connector.GetClient(reconnect)
	if c == nil {
		c = <-ch
	}
	client := c.(*gorm.DB)
	return client
}

func GetClient(dialect, params string) (*gorm.DB, error) {
	db, err := gorm.Open(dialect, params)
	if err != nil {
		return db, err
	}

	db.DB().SetMaxIdleConns(10)
	db.DB().SetMaxOpenConns(100)
	db.SingularTable(true)
	db.LogMode(false)
	return db, err
}

/**

orm, err := gorm2.CreateGorm("postgres", Model.Def.CommConf.HBPg.String)
if err != nil {
	log2.Panic(err)
}

type ShortcodeMediaVO struct {
	ID int `json:"id"`
}
var s ShortcodeMediaVO
ss := []ShortcodeMediaVO{}
db := orm.GetClient(false)
db.Raw("select * from shortcode_media_vo where id > 100 limit 10").Scan(&ss)
*/
func CreateGorm(dialect, params string) (*Gorm, error) {
	p := &Gorm{
		Connector: CreateConnector(dialect, params),
	}
	return p, nil
}

func CreateConnector(dialect, params string) *io2.Connector {
	alarmed := false
	return &io2.Connector{
		Connect: func() (interface{}, error) {
			return GetClient(dialect, params)
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
