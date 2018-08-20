package elasticsearch

import (
	"time"

	"github.com/olivere/elastic"
	"shangwoa.com/io2"
)

func GetClient(urls ...string) (*elastic.Client, error) {
	client, err := elastic.NewClient(elastic.SetURL(urls...), elastic.SetSniff(false))
	return client, err
}

type ElasticRemote struct {
	Connector *io2.Connector
}

func (this *ElasticRemote) GetClient(reconnect bool) *elastic.Client {
	c, ch := this.Connector.GetClient(reconnect)
	if c == nil {
		c = <-ch
	}
	client := c.(*elastic.Client)
	return client
}

/**
es, err := elasticsearch.CreateEs(Model.Def.CommConf.HBElastic.String)
	if err != nil {
		log2.Panic(err)
	}
	Model.Def.CommRemoteCollection.HbElastic = es
	es.GetClient(false).Index().Index("ig_post").Type("ig_post")
	res, err := es.GetClient(false).Search().Index("ig_post").From(0).Size(1).Pretty(true).Do(context.Background())
	if err != nil {
		log2.Panic(err)
	}
	fmt.Println(res)
*/
func CreateEs(params string) (*ElasticRemote, error) {
	p := &ElasticRemote{
		Connector: CreateConnector(params),
	}
	return p, nil
}
func CreateConnector(strs ...string) *io2.Connector {
	alarmed := false
	connector := &io2.Connector{
		Connect: func() (interface{}, error) {
			return GetClient(strs...)
		}, OnError: func(e error) bool {
			return false
		}, OnMaxSeqAlarm: func(i int) {
			if !alarmed {
				alarmed = true
				// alarm
			}
		}, MaxSeq: 2000, RetryDelay: time.Second * 2,
	}
	return connector
}
