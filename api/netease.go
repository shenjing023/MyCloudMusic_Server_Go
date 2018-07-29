/*
 * 网易云音乐api
 */

package api

import (
	"MyCloudMusic_Server_Go/mylog"
	"MyCloudMusic_Server_Go/utils"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/url"
	"strconv"
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
			//"User-Agent":       "Mozilla/5.0 (Windows NT 10.0; WOW64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/58.0.3029.110 Safari/537.36",
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
		params, encSecKey, err := encryptoParams(param)
		if err != nil {
			mylog.Error(err.Error())
		}

		form := url.Values{}
		form.Set("params", params)
		form.Set("encSecKey", encSecKey)
		body := strings.NewReader(form.Encode())
		req, _ = http.NewRequest("POST", _url, body)
	}

	for k, v := range n.header {
		req.Header.Set(k, v)
	}
	req.Header.Set("User-Agent", utils.GetUserAgent())

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
		return nil, err
	}
	defer resp.Body.Close()

	respBody, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		mylog.Error(err.Error())
		return nil, err
	}

	return respBody, nil
}

/*
	搜索单曲(1)，歌手(100)，专辑(10)，歌单(1000)，用户(1002) *(type)*
*/
func (n *NetEaseApi) Search(keywords string, ktype, offset, limit int) ([]byte, error) {
	url := "http://music.163.com/weapi/search/get"
	var total string
	if offset == 0 {
		total = "true"
	} else {
		total = "false"
	}
	params := map[string]interface{}{
		"s":      keywords,
		"type":   ktype,
		"offset": offset,
		"total":  total,
		"limit":  limit,
	}
	return n.rawHttpRequest("POST", url, params)
}

/*
	歌单（网友精选碟）
*/
func (n *NetEaseApi) Playlists(category string, order string, offset int, limit int) ([]byte, error) {
	url := "http://music.163.com/weapi/playlist/list"
	params := map[string]interface{}{
		"cat":    category,
		"order":  order,
		"offset": offset,
		"total":  "true",
		"limit":  limit,
	}
	respByte, err := n.rawHttpRequest("POST", url, params)
	if err != nil {
		return nil, err
	}

	var resp map[string]interface{}
	var playlists []map[string]interface{}
	if err = json.Unmarshal(respByte, &resp); err != nil {
		return nil, err
	}
	// 解析
	for _, item := range resp["playlists"].([]interface{}) {
		dict := item.(map[string]interface{})
		playlist := map[string]interface{}{
			"collect_name": dict["name"].(string),
			"list_id":      dict["id"].(float64),
			"logo":         dict["coverImgUrl"].(string),
		}
		playlists = append(playlists, playlist)
	}
	return json.Marshal(playlists)
}

/*
	根据歌单id获取歌单详情，返回歌单名称，歌单图片url，歌单作者，歌曲总数，歌单播放数，分享总数，
	收藏总数，还有每首歌的歌名，歌手，专辑，音乐id
	使用新版本v3接口，
*/
func (n *NetEaseApi) PlaylistDetail(playlistId string) ([]byte, error) {
	url := "http://music.163.com/weapi/v3/playlist/detail"
	params := map[string]interface{}{
		"id":         playlistId,
		"total":      "true",
		"limit":      1000,
		"n":          1000,
		"offset":     0,
		"csrf_token": "",
	}
	respByte, err := n.rawHttpRequest("POST", url, params)
	if err != nil {
		return nil, err
	}

	var resp map[string]interface{}
	if err = json.Unmarshal(respByte, &resp); err != nil {
		return nil, err
	}
	// 解析
	playlist := resp["playlist"].(map[string]interface{})
	var songs []map[string]interface{} //每首歌信息

	for _, item := range resp["playlist"].(map[string]interface{})["tracks"].([]interface{}) {
		dict := item.(map[string]interface{})
		songName := dict["name"] //音乐名称
		// 歌手
		var singers string
		for _, singer := range dict["ar"].([]interface{}) {
			singers = singer.(map[string]interface{})["name"].(string)
			singers += "/"
		}
		singers = singers[:len(singers)-1]
		// 专辑、图片url
		album := dict["al"].(map[string]interface{})["name"].(string)
		picUrl := dict["al"].(map[string]interface{})["picUrl"].(string)
		// 歌曲id
		songId := dict["id"].(float64)
		// 歌曲长度
		songLength := int(dict["dt"].(float64) / 1000)

		song := map[string]interface{}{
			"song_name":   songName,
			"singers":     singers,
			"album_name":  album,
			"song_id":     songId,
			"pic_url":     picUrl,
			"song_length": songLength,
		}
		songs = append(songs, song)
	}
	data := map[string]interface{}{
		"type":           "netease",
		"name":           playlist["name"].(string),
		"pic":            playlist["coverImgUrl"].(string),
		"user":           playlist["creator"].(map[string]interface{})["nickname"].(string),
		"songs_count":    playlist["trackCount"].(float64),
		"play_count":     playlist["playCount"].(float64),
		"share_count":    playlist["shareCount"].(float64),
		"favorite_count": playlist["subscribedCount"].(float64),
		"songs":          songs,
	}
	return json.Marshal(data)
}

/*
	根据歌曲id获取歌曲url
*/
func (n *NetEaseApi) SongUrl(songId string) ([]byte, error) {
	url := "http://music.163.com/weapi/song/enhance/player/url"
	params := map[string]interface{}{
		"ids":        []string{songId},
		"br":         320000,
		"csrf_token": "",
	}

	respByte, err := n.rawHttpRequest("POST", url, params)
	if err != nil {
		return nil, err
	}

	var resp map[string]interface{}
	if err = json.Unmarshal(respByte, &resp); err != nil {
		return nil, err
	}

	//解析
	data := map[string]string{
		"song_url": resp["data"].([]interface{})[0].(map[string]interface{})["url"].(string),
	}
	return json.Marshal(data)
}

/*
	专辑
*/
func (n *NetEaseApi) Album(albumId string) ([]byte, error) {
	url := fmt.Sprintf("http://music.163.com/weapi/v1/album/%s", albumId)
	params := map[string]interface{}{
		"id":         albumId,
		"csrf_token": "",
	}
	return n.rawHttpRequest("POST", url, params)
}

/*
	歌词
*/
func (n *NetEaseApi) Lyric(songId string) ([]byte, error) {
	url := "http://music.163.com/weapi/song/lyric"
	params := map[string]interface{}{
		"id": songId,
		"os": "osx",
		"lv": -1,
		"kv": -1,
		"tv": -1,
	}
	respByte, err := n.rawHttpRequest("POST", url, params)
	if err != nil {
		return nil, err
	}

	var resp map[string]interface{}
	if err = json.Unmarshal(respByte, &resp); err != nil {
		return nil, err
	}

	// 解析
	var lrc interface{}
	if resp["lrc"] == nil { //没有歌词
		lrc = ""
	} else {
		lrc = resp["lrc"].(map[string]interface{})["lyric"]
	}
	data := map[string]interface{}{
		"lyric": lrc,
	}

	return json.Marshal(data)
}

/*
	私人FM
*/
func (n *NetEaseApi) PersonFM() ([]byte, error) {
	url := "http://music.163.com/weapi/v1/radio/get"
	params := map[string]interface{}{
		"csrf_token": "",
	}

	respByte, err := n.rawHttpRequest("POST", url, params)
	if err != nil {
		fmt.Println(err.Error())
		return nil, err
	}

	var resp map[string]interface{}
	if err = json.Unmarshal(respByte, &resp); err != nil {
		return nil, err
	}

	// 解析
	// 歌曲信息
	if resp["data"] == nil {
		// 休眠1秒
		time.Sleep(time.Duration(1) * time.Second)
		mylog.Info("私人fm: 没有歌曲信息")
		return n.PersonFM()
	}
	var songInfo = resp["data"].([]interface{})[0].(map[string]interface{})
	// 专辑图片url
	picUrl := songInfo["album"].(map[string]interface{})["picUrl"]
	// 专辑名
	albumName := songInfo["album"].(map[string]interface{})["name"]
	// 歌手
	var singer []string
	for _, item := range songInfo["artists"].([]interface{}) {
		singer = append(singer, item.(map[string]interface{})["name"].(string))
	}
	data := map[string]interface{}{
		"name":     songInfo["name"],
		"duration": int(songInfo["duration"].(float64) / 1000),
		"album":    albumName,
		"picUrl":   picUrl,
		"singer":   strings.Join(singer, ","),
		"lyric":    "",
		"id":       songInfo["id"],
	}
	// 获取歌词
	lyric, err := n.Lyric(strconv.Itoa(int(songInfo["id"].(float64))))
	if err != nil {

	} else {
		var lrc map[string]interface{}
		json.Unmarshal(lyric, &lrc)
		data["lyric"] = lrc["lyric"].(string)
	}
	// 获取歌曲url
	//songUrl, err := n.SongUrl(strconv.Itoa(int(songInfo["id"].(float64))))
	//if err != nil {
	//
	//} else {
	//	var url map[string]interface{}
	//	json.Unmarshal(songUrl, &url)
	//	data["url"] = url["song_url"].(string)
	//}

	return json.Marshal(data)
}
