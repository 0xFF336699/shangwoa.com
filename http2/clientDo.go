package http2

import (
	"net/http"
	"io/ioutil"
	"net/http/cookiejar"
)

func GetJarClient() *http.Client {
	Jar, _ := cookiejar.New(nil)
	client := &http.Client{Jar: Jar}
	return client
}

func GetEmptyClient() *http.Client {
	return &http.Client{}
}

func ClientDo(r *http.Request, client *http.Client) (content []byte, err error) {
	res, err := client.Do(r)
	if err != nil {
		return nil, err
	}
	defer res.Body.Close()
	content, err = ioutil.ReadAll(res.Body)
	if err != nil {
		return nil, err
	}
	return
}
