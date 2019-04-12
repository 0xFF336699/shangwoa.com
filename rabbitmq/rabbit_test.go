package rabbitmq

import (
	"testing"
	"time"
	"fmt"
	"shangwoa.com/utils/retry"
)

func TestNewRetries(t *testing.T) {
	//conn, err := GetConn("amqp://ig-crawler:ig-crawler@rabbitmq.hb.ms.shangwoa.com:8231/ig-crawler", retry.NewRetries(2, 1 * time.Second, 2, 2 * time.Second))
	conn, err := GetConn("amqp://ig-crawler:ig-crawler@rabbitmq.hb.ms.shangwoa.com:8231/ig-crawler", retry.NewRetries(2, 1 * time.Second, 2, 2 * time.Second))
	if err != nil{
		fmt.Println("err is ", err)
	}else{
		fmt.Println("conn", conn)
	}
}
