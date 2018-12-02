package qqbot

import (
	"encoding/json"
	"io/ioutil"
	"net/http"
	"strconv"
)

const (
	PostTypeMessage = "message" //消息
	PostTypeEvent = "event" //
	PostTypeRequest = "request" // 加好友/群请求
	MessageTypePrivate = "private"
	MessageTypeGroup = "group"
	MessageTypeDiscuss = "discuss"
	MessageType = "message_type"
	ReturnStatusNoBot = "no bot"
	ReturnStatusOK = "ok"
)
type GroupMessage struct {
	Anonymous   interface{} `json:"anonymous"`
	Font        int         `json:"font"`
	GroupID     int         `json:"group_id"`
	Message     string      `json:"message"`
	MessageID   int         `json:"message_id"`
	MessageType string      `json:"message_type"`
	PostType    string      `json:"post_type"`
	RawMessage  string      `json:"raw_message"`
	SelfID      int         `json:"self_id"`
	Sender      struct {
		Age      int    `json:"age"`
		Card     string `json:"card"`
		Nickname string `json:"nickname"`
		Sex      string `json:"sex"`
		UserID   int    `json:"user_id"`
	} `json:"sender"`
	SubType string `json:"sub_type"`
	Time    int    `json:"time"`
	UserID  int    `json:"user_id"`
}
type PrivateMessageListener func(i interface{}) (err error, status string, m map[string]interface{})
type BotListener struct{
	Bot Bot
	ID int64
	PrivateMessageListener PrivateMessageListener
}
type Listener struct{
	port int
	path string
	bots map[int64] *BotListener
}
var listeners = make(map[string]*Listener)
func CreateListener(port int, path string, bot *BotListener) (*Listener) {
	url := strconv.Itoa(port) + path
	if l, ok := listeners[url]; ok{
		l.bots[bot.ID] = bot
		return l
	}
	bots := map[int64]*BotListener{}
	bots[bot.ID] = bot
	l := &Listener{port:port, path:path, bots: bots}
	listeners[url] = l
	http.HandleFunc(path, func(writer http.ResponseWriter, request *http.Request) {
		defer request.Body.Close()
		buffer, _ := ioutil.ReadAll(request.Body)
		println("buffer is", string(buffer))
		m := map[string]interface{}{}
		json.Unmarshal(buffer, &m)
		rm := map[string]interface{}{}
		var status string
		var err error
		switch m["post_type"] {
		case "message":
			err, status, rm = l.onMessage(m)
			if err != nil{
				break;
			}
			if b, ok := rm["stop"]; ok && b.(bool) {
				break
			}
			//cs := command.Excision(m["message"].(string))
			//go command.Exec(cs[0], cs, m)
		case "event":
			//rm = s.listener.onEvent(m)
			break
		case "request":
			//rm = s.listener.onRequest(m)
			break
		}
		if status != ReturnStatusOK{
			// neet to send it to the hendler
		}
		if rm == nil || len(rm) < 1 {
			writer.WriteHeader(204)
		} else {
			jo, _ := json.Marshal(rm)
			writer.Write(jo)
		}
	})

	go func() {
		err := http.ListenAndServe(":"+strconv.Itoa(port), nil)
		if err != nil {
			panic(err)
		}
	}()
	return l
}

func (l *Listener) onMessage(m map[string]interface{}) (err error, status string, res map[string]interface{}) {
	selfID := int64(m["self_id"].(float64))
	bot, ok := l.bots[selfID]
	if ok == false{
		return nil, ReturnStatusNoBot, nil
		// it's an unexpected, we need to handle it
	}
	switch m[MessageType] {
	case MessageTypePrivate:
		if bot.PrivateMessageListener != nil{
			return bot.PrivateMessageListener(m)
		}
		break
	case MessageTypeGroup:

		break
	case MessageTypeDiscuss:

		break
	}
	return
}