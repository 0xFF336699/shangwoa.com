一般使用方法  
本地有个local.conf.json的文件  
```
// 这里储存了所有consul里的key路径
type ConsulKeys struct{
// consul里面kv的k路径
	ConsulKeyPath                       string `json:"consul_key_path"`
}
// 这里储存了所有consul的数据
type ConsulKVs struct{
//consul里面储存的数据
	ConsulValue *consul.KVString
}

//配置文件里对consulKeys kvs储存引用
type Conf struct {
	ConsulKeys *ConsulKeys `json:"consul_keys"`
	KVs *ConsulKVs
	KVsMap map[string]*consul.KV
}
// 单例里对consul数据进行音乐
var MImpl *Model = &Model{Conf:&Conf{
	KVsMap: map[string]*consul.KV{},
	KVs:&ConsulKVs{
// 把ConsulKVs里的数据对象先进行初始化
		ConsulValue: &consul.KVString{},
}},
}

//链接consul进行初始化
	conf := model.Conf
	keys := conf.ConsulKeys
	kvs := conf.KVs
	m := conf.KVsMap
	m[keys.ConsulKeyPath] = &consul.KV{KVValue: kvs.ConsulValue}

```