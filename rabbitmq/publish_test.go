package rabbitmq

import (
	"testing"
	"fmt"
)

func TestPublishByDefault(t *testing.T) {
	body := []byte("hehe")
	err := PublishByDefault("post_media_order:downloaded", "amqp://ig-crawler:ig-crawler@rabbitmq.hb.ms.shangwoa.com:8231/ig-crawler", body)
	if err != nil{
		fmt.Println(err)
	}
}
