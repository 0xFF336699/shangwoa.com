package consul

import (
	"testing"
	"fmt"
	"shangwoa.com/json2"
	"net/http"
	"io/ioutil"
	"bytes"
)

type kvData struct {
	Name string
	Seq int `json:"seq"`
}
var token = "";
func TestReadKeys(t *testing.T) {
	//testPut()
	//return
	path := "http://consul.ms.shangwoa.com/v1/txn?token="

	kv := &kvData{}
	kvs := [] *KV{&KV{Key:"ig/redis/prefix/ig_user", KVValue:&KVValue{}},
		&KV{Key:"ig/redis/prefix/post", KVJSONB:&json2.JSONB{}},
		&KV{Key:"ig/service/redis/oversea/ig", KVStruct:&KVRedis{}},
		&KV{Key:"ig/redis/prefix/seq", KVStruct:&kv},
		&KV{Key:"ig/redis/prefix/prefix"}}
	url := path + token
	err := GetKeys(url, kvs)
	if err != nil{
		fmt.Println(err)
	}
	fmt.Println(kvs)

}

func testPut() {
	url := "http://consul.ms.shangwoa.com/v1/txn?token=" + token
	fmt.Println("URL:>", url)

	var jsonStr = []byte(`[
  {
    "KV": {
      "Verb": "get",
      "Key": "ig/service/redis/oversea/ig"
    }
  },
  {
    "KV": {
      "Verb": "get",
      "Key": "ig/service/rabbitmq/oversea"
    }
  }
]`)
	req, err := http.NewRequest("PUT", url, bytes.NewBuffer(jsonStr))
	req.Header.Set("X-Custom-Header", "myvalue")
	req.Header.Set("Content-Type", "application/json")

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		panic(err)
	}
	defer resp.Body.Close()

	fmt.Println("response Status:", resp.Status)
	fmt.Println("response Headers:", resp.Header)
	body, _ := ioutil.ReadAll(resp.Body)
	fmt.Println("response Body:", string(body))
}

//[
//{
//	"KV": {
//		"Verb": "get",
//		"Key": "ig/service/redis/oversea/ig"
//	}
//},
//{
//	"KV": {
//"Verb": "get",
//"Key": "ig/service/rabbitmq/oversea"
//}
//}
//]