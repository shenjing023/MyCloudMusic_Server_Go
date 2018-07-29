package main

import (
	"MyCloudMusic_Server_Go/mylog"
	"github.com/julienschmidt/httprouter"
)

/*
路由注册
*/

type Route struct {
	Name        string
	Method      string
	Path        string
	HandlerFunc httprouter.Handle
}

type Routes []Route

func allRoutes() Routes {
	routes := Routes{
		Route{"Index", "GET", "/", Index},
		Route{"Search", "GET", "/search", Search},
		Route{"Playlist", "GET", "/playlists", PlayList},
		Route{"PlaylistDetail", "GET", "/playlist/detail", PlayListDetail},
		Route{"SongUrl", "GET", "/music/url", SongUrl},
		Route{"PersonFM", "GET", "/personfm", PersonFM},
		Route{"Lyric", "GET", "/lyric", Lyric},
	}
	return routes
}

// 返回一个路由
func NewRouter(routes Routes) *httprouter.Router {
	router := httprouter.New()
	for _, route := range routes {
		var handle httprouter.Handle
		handle = route.HandlerFunc
		handle = mylog.Log(handle)
		router.Handle(route.Method, route.Path, handle)
	}
	return router
}
