/*
 * 网易云音乐api
 */

package api

import (
	"MyCloudMusic_Server_Go/mylog"
	"io/ioutil"
	"net/http"
	"net/url"
	"strings"
	"time"
)

type NetEaseApi struct {
	header map[string]string //请求头
	client *http.Client      //请求客户端
}

var NetEase *NetEaseApi

func init() {
	NetEase = &NetEaseApi{
		header: map[string]string{
			"Accept":           "*/*",
			"Accept-Langeuage": "zh-CN,zh;q=0.8,gl;q=0.6,zh-TW;q=0.4",
			"Connection":       "keep-alive",
			"Content-Type":     "application/x-www-form-urlencoded",
			"Host":             "music.163.com",
			"Referer":          "http://music.163.com",
			"User-Agent":       "Mozilla/5.0 (Windows NT 10.0; WOW64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/58.0.3029.110 Safari/537.36",
		},
		client: &http.Client{
			Timeout: 15 * time.Second,
		},
	}
}

/*
	设置http请求体，返回http.Request
*/
func (n *NetEaseApi) newRequest(method, _url string, param map[string]interface{}) (*http.Request, error) {
	var req *http.Request

	if method == "GET" {
		req, _ = http.NewRequest("GET", _url, nil)
	} else if method == "POST" {
		params, encSeckey, err := encryptoParams(param)
		if err != nil {
			mylog.Error(err.Error())
		}
		form := url.Values{}
		form.Set("params", params)
		form.Set("encSeckey", encSeckey)
		body := strings.NewReader(form.Encode())
		req, _ = http.NewRequest("POST", _url, body)
	}

	for k, v := range n.header {
		req.Header.Set(k, v)
	}

	return req, nil
}

func (n *NetEaseApi) httpRequest(method, url, query string) {

}

func (n *NetEaseApi) rawHttpRequest(method, url string, param map[string]interface{}) ([]byte, error) {
	var req *http.Request
	req, _ = n.newRequest(method, url, param)

	resp, err := n.client.Do(req)
	if err != nil {
		mylog.Error(err.Error())
	}
	defer resp.Body.Close()

	return ioutil.ReadAll(resp.Body)
}

/*
	搜索单曲(1)，歌手(100)，专辑(10)，歌单(1000)，用户(1002) *(type)*
*/
func (n *NetEaseApi) Search(keywords string, ktype, offset, limit int) ([]byte, error) {
	url := "http://music.163.com/api/search/get"
	var total bool
	if offset == 0 {
		total = true
	} else {
		total = false
	}
	data := map[string]interface{}{
		"s":      keywords,
		"type":   ktype,
		"offset": offset,
		"total":  total,
		"limit":  limit,
	}
	return n.rawHttpRequest("POST", url, data)
}
