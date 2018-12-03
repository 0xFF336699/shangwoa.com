package qqbot_mq

import (
	"flag"
	"fmt"
	"path/filepath"
	"shangwoa.com/os2"
)

func init()  {
	initFlags()
	loadLocalConfig()
	loadMqConf()
}


var flagConfFolder string
func initFlags() {
	flag.StringVar(&flagConfFolder, "conf", "./conf", "config path, must be a folder.")
	flag.Parse()
}

func loadLocalConfig() {
	files, err := filepath.Glob(flagConfFolder + "/host/*.json")
	if err != nil{
		panic(err)
	}

	if len(files) == 0{
		panic("no host")
	}
	hosts := []*Host{}
	for _, p := range files{
		hs := []*Host{}
		err := os2.LoadFileToStruct(p, &hs)
		if err != nil{
			panic(err)
		}
		hosts = append(hosts, hs...)
	}
	for _, h := range hosts{
		CreateListener(h.Port, h.Path)
	}
	println(hosts)
}
func loadMqConf() {
	files, err := filepath.Glob(flagConfFolder + "/mq/*.json")
	if err != nil{
		panic(err)
	}

	if len(files) == 0{
		panic("no host")
	}
	mqconfs := []*Mq{}
	for _, p := range files{
		temps := []*Mq{}
		err := os2.LoadFileToStruct(p, &temps)
		if err != nil{
			panic(err)
		}
		mqconfs = append(mqconfs, temps...)
	}
	for _, m := range mqconfs{
		fmt.Println("%#v", m)
		for _, b := range m.Bots{
			for _, f := range b.Friend.Excluede{
				b.Friend.ExcludeMap = map[int]*FriendExcluede{}
				b.Friend.ExcludeMap[f.ID] = f
			}
			for _, s := range b.Friend.Special{
				b.Friend.SpecialMap = map[int]*FriendSpecial{}
				b.Friend.SpecialMap[s.ID] = s
			}
			for _, e := range b.Group.Exclude{
				b.Group.ExcluedMap = map[int]*GroupExclude{}
				b.Group.ExcluedMap[e.ID] = e
			}
			for _, s := range b.Group.Special{
				b.Group.SpecialMap = map[int]*GroupSpecial{}
				b.Group.SpecialMap[s.ID] = s
			}
			b.Mq = m
			model.Bots[b.BotID] = b
		}
		//CreateListener(h.Port, h.Path)
	}
	b:= model.Bots
	println(b)
}
