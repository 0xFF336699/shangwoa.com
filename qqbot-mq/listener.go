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
	fmt.Println("create listener", port, path);
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
		//println("buffer is", string(buffer))
		m := map[string]interface{}{}

		if err := json.Unmarshal(buffer, &m);err != nil{
			return
		}
		selfID := int(m["self_id"].(float64))
		bot := model.Bots[selfID]
		qname := strconv.Itoa(selfID)
		jump := ""
		const c = ConnectCharacter
		switch m[PostType] {
		case PostTypeMessage:
			qname += c + PostTypeMessage
			switch m[MessageType] {
			case MessageTypePrivate:
				qname += c + MessageTypePrivate
				userID := int(m["user_id"].(float64))
				if _, ok := bot.Friend.ExcludeMap[userID]; ok{
					jump = "exclude"
					break
				}
				if _, ok := bot.Friend.SpecialMap[userID]; ok{
					qname += c + strconv.Itoa(userID)
				}
				break
			case MessageTypeGroup:
				qname += c + MessageTypeGroup
				groupID := int(m["group_id"].(float64))
				if _, ok := bot.Group.ExcluedMap[groupID]; ok{
					jump = "exclude"
					break
				}
				if _, ok := bot.Group.SpecialMap[groupID]; ok{
					//fmt.Println("it's sticker", string(buffer))
					qname += c + strconv.Itoa(groupID)
				}else{
					//fmt.Println("it's normal", string(buffer))
				}
				break
			case MessageTypeDiscuss:
				qname += c + MessageTypeDiscuss
				break
			default:
				jump = MessageType + "error"
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
		if jump != ""{
			fmt.Println("jump is ", jump)
			return
		}
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
var pcs map[string]*rabbitmq.PubChannel
func getChannel(url, qname string) (err error, pc *rabbitmq.PubChannel){
	name := url+"_" + qname
	if pc, ok := pcs[name]; ok{
		return nil, pc
	}
	args := &rabbitmq.ChannelArgs{
		URL:url,
		Exchange:"amq.direct",
		Qname:"",
		Kind:"direct",
		RoutingKey:qname,
		Durable:true,
		Mandatory:false,
		Immediate:false,
		AutoDelete:false,
		Internal:false,
		Exclusive:false,
		NoWait:false,
		RetryMaxCount:10,
		InteralTime:1000,
		Args:nil,
		DeliveryMode:2,
		ContentType:"text/plain",
		MaxWaitingCount:10,
	}
	err, pc = rabbitmq.MaxWaitingPublishChannel(args)
	if err!= nil{
		return
	}
	pcs[name] = pc
	return
}
func SendMsg(mq, qname string, body []byte) {
	sendMsg(mq, qname, body)
}
func sendMsg(mq string, qname string, body []byte)  {
	err, pc := getChannel(mq, qname)
	if err != nil{
		fmt.Println("publish apply channel error", err.Error(), string(body))
		return
	}
	err = pc.Publish(body)
	if err != nil{
		fmt.Println("publish error", err.Error(), string(body))
		cancelPc(mq, qname)// 没有做后续处理
	}else{
		fmt.Println("2019-04-19 publish ok", string(body))
	}
}

func cancelPc(mq, qname string)  {
	delete(pcs, mq + "_" + qname)
}
func sendMsg2(mq string, qname string, body []byte)  {

	//fmt.Println("qname is", qname, mq)
	//err := rabbitmq.PublishByDefault("post_media_order:downloaded", "amqp://ig-crawler:ig-crawler@rabbitmq.hb.ms.shangwoa.com:8231/ig-crawler", body)
	//err := rabbitmq.PublishByDefault(qname, mq, body)
	err := rabbitmq.Publish(qname, "amq.direct", qname, "direct", mq, body)
	if err != nil {
		fmt.Println("publish error", err.Error(), string(body))
	}else{
		fmt.Println("publish ok", string(body))
	}
}