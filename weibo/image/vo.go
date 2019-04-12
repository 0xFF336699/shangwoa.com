package image

import (
	"encoding/json"
	"errors"
)

type ImageInfo struct {
	Src string `json:"src"`
	Pid string `json:"pid,omitempty"`
	Err error `json:"-"`
	ErrMsg string `json:"err_msg,omitempty"`
	Width int `json:"width,omitempty"`
	Height int `json:"height,omitempty"`
}

type MqInfo struct {
	Qname string `json:"qname,omitempty"`
	RoutingKey string `json:"routing_key,omitempty"`
	Kind string `json:"kind,omitempty"`
	Exchange string `json:"exchange,omitempty"`
	AutoDelete bool `json:"auto_delete,omitempty"`
	Url string `json:"url,omitempty"`
}
type DownUploadInfo struct{
	List []*ImageInfo `json:"list,omitempty"`
	MqInfo *MqInfo    `json:"mq_info,omitempty"`
	UseProxy bool     `json:"use_proxy,omitempty"`
	ID int `json:"id,omitempty"` // 用来识别校验的

}
func ParseImageList(josnBytes []byte)(err error, info *DownUploadInfo){
	err = json.Unmarshal(josnBytes, &info)
	if err != nil{
		return
	}
	if len(info.MqInfo.Qname) == 0{
		err = errors.New("no qname")
		return
	}
	//if len(info.MqInfo.RoutingKey) == 0{
	//	err = errors.New("no routing key")
	//	return
	//}
	if len(info.MqInfo.Kind) == 0{
		info.MqInfo.Kind = "direct"
	}
	if len(info.MqInfo.Exchange) == 0{
		//info.MqInfo.Exchange = "amq.direct"
		err = errors.New("no exchange")
		return
	}
	return
}