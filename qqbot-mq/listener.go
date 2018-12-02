package qqbot_mq

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"shangwoa.com/rabbitmq"
	"strconv"
	. "shangwoa.com/qqbot-const"
)

type BotListener struct{
	ID int64
}
type Listener struct{
	port int
	path string
	bots map[int64] *BotListener
}
var listeners = make(map[string]*Listener)
func CreateListener(port int, path string) (*Listener) {
	url := strconv.Itoa(port) + path
	if l, ok := listeners[url]; ok{
		return l
	}
	bots := map[int64]*BotListener{}
	l := &Listener{port:port, path:path, bots: bots}
	listeners[url] = l
	http.HandleFunc(path, func(writer http.ResponseWriter, request *http.Request) {
		defer request.Body.Close()
		writer.WriteHeader(204)
		buffer, _ := ioutil.ReadAll(request.Body)
		println("buffer is", string(buffer))
		m := map[string]interface{}{}

		if err := json.Unmarshal(buffer, &m);err != nil{
			return
		}
		selfID := int(m["self_id"].(float64))
		qname := strconv.Itoa(selfID)
		const c = "._."
		switch m[PostType] {
		case PostTypeMessage:
			qname += c + PostTypeMessage
			switch m[MessageType] {
			case MessageTypePrivate:
				qname += c + MessageTypePrivate
				break
			case MessageTypeGroup:
				qname += c + MessageTypeGroup
				break
			case MessageTypeDiscuss:
				qname += c + MessageTypeDiscuss
				break
			default:
				return // error
			}
			break
		case PostTypeEvent:
			qname += c + PostTypeEvent
			switch m[EventEvent] {
			case EventGroupAdmin:
				qname += c + EventGroupAdmin
				break
			case EventGroupDecrease:
				qname += c + EventGroupDecrease
				break
			case EventGroupIncrease:
				qname += c + EventGroupIncrease
				break
			case EventGroupUpload:
				qname += c + EventGroupUpload
				break
			default:
				return // error
			}
			break
		case PostTypeRequest:
			qname += c + PostTypeRequest
			switch m[RequestType] {
			case RequestTypeFriend:
				qname += c + RequestTypeFriend
				break
			case RequestTypeGroup:
				qname += c + RequestTypeGroup
				break
			default:
				return // error
			}
			break
		default:
			return // it's an error
		}
		bot := model.Bots[selfID]
		sendMsg(bot.Mq.Amqp, qname, buffer);
	})

	go func() {
		err := http.ListenAndServe(":"+strconv.Itoa(port), nil)
		if err != nil {
			panic(err)
		}
	}()
	return l
}

func sendMsg(mq string, qname string, body []byte)  {

	fmt.Println("qname is", qname, mq)
	//err := rabbitmq.PublishByDefault("post_media_order:downloaded", "amqp://ig-crawler:ig-crawler@rabbitmq.hb.ms.shangwoa.com:8231/ig-crawler", body)
	err := rabbitmq.PublishByDefault(qname, mq, body)
	if err != nil {
		fmt.Println("publish error", err.Error())
	}
}