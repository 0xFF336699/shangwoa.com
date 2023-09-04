package net2
import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"runtime"


	"io"
	"mime/multipart"
	"os"
	"strings"
)

var client = &http.Client{}
func Request(params map[string]interface{}, endpoint string) (err error, body []byte) {
	bs, err := json.Marshal(params)
	if err != nil{
		return
	}
	req, err := http.NewRequest("POST", endpoint, bytes.NewBuffer(bs))
	req.Header.Set("Content-Type", "application/json;charset=UTF-8")
	resp, err := client.Do(req)
	if err != nil{
		___fn := ""
		__pc, _,___line, ___ok := runtime.Caller(0)
		if(___ok){
			___fn = runtime.FuncForPC(__pc).Name()
		}
		fmt.Println("error request", ___fn, ___line, "pdd request", err.Error())
		return
	}
	defer resp.Body.Close()

	body, err = ioutil.ReadAll(resp.Body)
	return
}
func RequestForm(params map[string]interface{}, endpoint string) (err error, body []byte) {
	ct, b, err := CreateForm(params)

	resp, err := http.Post(endpoint, ct, b)
	//req, err := http.NewRequest("POST", endpoint, bytes.NewBuffer(bs))
	//req.Header.Set("Content-Type", "application/json;charset=UTF-8")
	//resp, err := client.Do(req)
	if err != nil{
		___fn := ""
		__pc, _,___line, ___ok := runtime.Caller(0)
		if(___ok){
			___fn = runtime.FuncForPC(__pc).Name()
		}
		fmt.Println("error request", ___fn, ___line, "pdd request", err.Error())
		return
	}
	defer resp.Body.Close()

	body, err = ioutil.ReadAll(resp.Body)
	return
}


func CreateForm(form map[string]interface{}) (string, io.Reader, error) {
	body := new(bytes.Buffer)
	mp := multipart.NewWriter(body)
	defer mp.Close()
	for key, v := range form {
		val := fmt.Sprintf("%v", v)
		if strings.HasPrefix(val, "@") {
			val = val[1:]
			file, err := os.Open(val)
			if err != nil { return "", nil, err }
			defer file.Close()
			part, err := mp.CreateFormFile(key, val)
			if err != nil { return "", nil, err }
			io.Copy(part, file)
		} else {
			mp.WriteField(key, val)
		}
	}
	return mp.FormDataContentType(), body, nil
}

func HttpGet(url string) (err error, isOk bool, body []byte) {
	resp, err := http.Get(url)
	if err != nil {
		fmt.Println(err)
		return
	}
	defer resp.Body.Close()
	body, err = ioutil.ReadAll(resp.Body)
	fmt.Println(string(body))
	fmt.Println(resp.StatusCode)
	if resp.StatusCode == 200 {
		isOk = true
	}
	return
}