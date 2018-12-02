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
			b.Mq = m
			model.Bots[b.BotID] = b
		}
		//CreateListener(h.Port, h.Path)
	}
	println(model)
}