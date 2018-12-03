package qqbot

import (
	"encoding/json"
	"errors"
	"io/ioutil"
	"net/http"
	"strconv"
	"strings"
)
const (
	BotResOk = "ok"
	BotResFailed = "failed"
	BotResAsync = "async"
	ApiGetLoginInfo = "get_login_info"
	ApiSendPrivateMessage = "send_private_msg"
)
type SendMsg struct {
	Type string `json:"type"`
	Data map[string]interface{} `json:"data"`
}
type SendPrivateMsg struct {
	UserID int `json:"user_id"`
	Message []*SendMsg `json:"message"`
	AutoEscape bool `json:"auto_escape"`
}
type LoginUser struct{
	UserID int64  `json:"user_id"`
	Nickname string `json:"nickname"`
}
type User struct{
	Id int
}
type SendPrivateMsgRes struct {
	MessageID int32 `json:"message_id"`
}
type Bot struct {
	url string
	token string
	secret string
	port int
	User *User
	Listener *Listener
}

func CreateMsg(t string, data map[string]interface{}) (msg *SendMsg) {
	msg = &SendMsg{t, data}
	return
}
type BotRes struct {
	Status string `json:"status"` // ok, failed asynic
	Retcode int `json:"retcode"`
	Data map[string]interface{} `json:"data"`
}

func (bot *Bot) post(api string, data interface{}, resStruct interface{}) (err error, res *BotRes) {
	client := &http.Client{}
	jo, err := json.Marshal(data)
	if err != nil{
		return
	}
	path := bot.url + "/" + api
	request, err := http.NewRequest("POST", path, strings.NewReader(string(jo)))
	if err != nil {
		panic(err)
	}
	request.Header.Set("Authorization", "Token "+ bot.token)
	request.Header.Set("Content-Type", "application/json")
	rb, err := client.Do(request)
	if rb != nil{
		defer rb.Body.Close()
	}
	if err != nil {
		panic(err)
	}
	buffer, err := ioutil.ReadAll(rb.Body)
	res = &BotRes{}
	err = json.Unmarshal(buffer, res)
	if err != nil{
		return
	}
	switch res.Status {
	case BotResOk:
		if resStruct != nil{
			var resData map[string]interface{}
			resData = res.Data
			b, err := json.Marshal(resData)
			if err != nil{
				return err, nil
			}
			err = json.Unmarshal(b, resStruct)
			if err != nil{
				return err, nil
			}
		}
		break;
	case BotResAsync:

		break;
	case BotResFailed:
		err = errors.New(strconv.Itoa(res.Retcode))
		return
		break;
	default:
		err = errors.New(res.Status)
		return
	}
	return
}

func (bot *Bot) GetLoginInfo() (err error, loginUser *LoginUser) {
	loginUser = &LoginUser{}
	err, _ = bot.post(ApiGetLoginInfo, make(map[string]interface{}), loginUser)
	return
}
func (bot *Bot) SendPrivateMsg(msg *SendPrivateMsg) (err error) {
	res := &SendPrivateMsgRes{}
	err, _ = bot.post(ApiSendPrivateMessage, msg, res)
	return
}
var bots = make(map[int64]*Bot)
func CreateBot(id int64, token, url, path, screte string, port int,
	botListener *BotListener) (err error, bot *Bot) {
		botListener.ID = id
	bot = &Bot{url:url,token:token, secret:screte, port:port}
	if botListener != nil{
		bot.Listener = CreateListener(port, path, botListener)
	}
	err, loginUser := bot.GetLoginInfo()
	if err != nil{
		return
	}
	if loginUser.UserID != id{
		err = errors.New("wrong bot id")
		return
	}
	bots[id] = bot
	return
}

func GetOrCreateBot(id int64, token, url, path, screte string, port int,
	botListener *BotListener) (err error, bot *Bot)  {
	bot, ok := bots[id]
	if ok{
		return
	}
	err, bot = CreateBot(id, token, url, path, screte, port, botListener)
	return
}