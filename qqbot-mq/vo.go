package qqbot_mq


var model = &Model{Hosts: []*Host{}, Bots: map[int]*Bot{}}
type Host struct{
	Enabled bool `json:"enabled"`
	Comment string `json:"comment"`
	Path string `json:"path"`
	Port int `json:"port"`
}
type Model struct {
	Hosts []*Host
	Mqs []*Mq
	Bots map[int]*Bot
}

type Mq struct {
	Comment string `json:"comment"`
	Amqp    string `json:"amqp"`
	Bots     []*Bot `json:"bots"`
}


// ----------------
type FriendDef struct {
	Comment     string `json:"comment"`
	QueuePrefix string `json:"queue_prefix"`
	QueueSuffix string `json:"queue_suffix"`
}
type FriendSpecial struct {
	Comment     string `json:"comment"`
	ID          int    `json:"id"`
	Alias       string `json:"alias"`
	QueuePrefix string `json:"queue_prefix"`
	QueueSuffix string `json:"queue_suffix"`
}
type FriendExcluede struct {
	Comment string `json:"comment"`
	ID      int    `json:"id"`
	Alias   string `json:"alias"`
}

type Friend  struct {
	Comment string `json:"comment"`
	All     bool   `json:"all"`
	Excluede [] *FriendExcluede `json:"excluede"`
	Special [] *FriendSpecial `json:"special"`
	Def *FriendDef `json:"def"`
	ExcludeMap map[int]*FriendExcluede
	SpecialMap map[int]*FriendSpecial
}


type GroupDef struct {
	Comment     string `json:"comment"`
	QueuePrefix string `json:"queue_prefix"`
	QueueSuffix string `json:"queue_suffix"`
}
type GroupSpecial struct {
	Comment string `json:"comment"`
	ID      int    `json:"id"`
	Alias   string `json:"alias"`
}
type GroupExclude struct {
	Comment string `json:"comment"`
	ID      int    `json:"id"`
	Alias   string `json:"alias"`
}

type Group struct {
	Comment string          `json:"comment"`
	All     bool            `json:"all"`
	Exclude []*GroupExclude `json:"exclude"`
	Special []*GroupSpecial `json:"special"`
	Def     *GroupDef        `json:"def"`
	ExcluedMap map[int]*GroupExclude
	SpecialMap map[int]*GroupSpecial
}

type Bot struct {
	Comment string `json:"comment"`
	BotID   int    `json:"bot_id"`
	Friend  *Friend `json:"friend"`
	Group *Group `json:"group"`
	Mq *Mq
}