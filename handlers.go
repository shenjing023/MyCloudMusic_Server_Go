package main

import (
	"net/http"
	"github.com/julienschmidt/httprouter"
	"fmt"
	"strconv"
)

type Argument struct {
	name string
	required bool
	dataType  string
}

type Resource []Argument

func (r *Resource) addArgument(name string,required bool,dataType string)  {
	arg:=Argument{
		name:name,
		required:required,
		dataType:dataType,
	}
	*r=append(*r,arg)
}

var search Resource

func init()  {
	search.addArgument("",true,"")
}

func (arg *Argument)parse()  {

}

func parseArgs(args []Argument)  {

}


func Index(w http.ResponseWriter,r *http.Request,_ httprouter.Params){
	fmt.Fprint(w,"hello world")
}

//搜索
func Searc1h(w http.ResponseWriter,r *http.Request,_ httprouter.Params)  {
	queryValues:=r.URL.Query()
	//source
	source:=queryValues.Get("source")
	if source==""{
		fmt.Fprintf(w,"source 必须")
	} else {
		if s,err:=strconv.Atoi(source);err!=nil{
			fmt.Fprintf(w,"%s",err.Error())
		} else {
			fmt.Fprintf(w,"%d",s)
		}
	}
}