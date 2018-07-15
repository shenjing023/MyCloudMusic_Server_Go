package main

import (
	"MyCloudMusic_Server_Go/api"
	"errors"
	"fmt"
	"github.com/julienschmidt/httprouter"
	"net/http"
	"strconv"
)

type Argument struct {
	name     string
	dataType string
}

type Resource []Argument

// 添加请求的参数
func (r *Resource) addArgument(name string, dataType string) error {
	// 参数的类型只能是string或int
	if dataType != "string" && dataType != "int" {
		return errors.New("dataType must be int or string")
	}
	// 参数名称不能为空
	if name == "" {
		return errors.New("name cannot be empty")
	}
	arg := Argument{
		name:     name,
		dataType: dataType,
	}
	*r = append(*r, arg)
	return nil
}

// 解析参数
func (r *Resource) parseArgs(values map[string][]string) []string {
	var errorMsg []string
	for _, arg := range *r {
		// 查看参数是否存在
		if _, ok := values[arg.name]; !ok {
			// 不存在
			errorMsg = append(errorMsg, "Missing required parameter "+arg.name)
		} else {
			if msg, err := arg.parse(values[arg.name]); err != nil {
				errorMsg = append(errorMsg, msg)
			}
		}
	}
	return errorMsg
}

var search Resource

func init() {
	search.addArgument("source", "string")
	search.addArgument("keywords", "string")
	search.addArgument("ktype", "int")
	search.addArgument("offset", "int")
	search.addArgument("limit", "int")
}

func (arg *Argument) parse(values []string) (string, error) {
	if arg.dataType == "int" {
		if _, err := strconv.Atoi(values[0]); err != nil {
			return values[0] + " cannot convert to int", err
		}
	}
	return "", nil
}

func Index(w http.ResponseWriter, _ *http.Request, _ httprouter.Params) {
	fmt.Fprint(w, "hello world")
}

// 搜索
func Search(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
	if errorMsg := search.parseArgs(r.URL.Query()); errorMsg != nil {
		fmt.Fprint(w, errorMsg)
	} else {
		queryValues := r.URL.Query()
		source := queryValues.Get("source")
		keywords := queryValues.Get("keywords")
		ktype, _ := strconv.Atoi(queryValues.Get("ktype"))
		offset, _ := strconv.Atoi(queryValues.Get("offset"))
		limit, _ := strconv.Atoi(queryValues.Get("limit"))

		var response []byte
		if source == "netease" {
			response, _ = api.NetEase.Search(keywords, ktype, offset, limit)
		} else {

		}

		fmt.Fprint(w, string(response))
	}
}
