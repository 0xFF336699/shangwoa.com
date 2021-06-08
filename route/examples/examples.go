package examples


import (
	"fmt"
	"net/http"
	"regexp"
	"shangwoa.com/route"
	"strings"
)

func StartServer(){

	var a = route.NewApp()
	r := &route.Route{
		Types: []string{http.MethodGet},
		Match: func(r *http.Request, data *route.RouterData) bool {
			b, err := regexp.MatchString(`^/user/([a-zA-Z0-9]+)$`, r.URL.Path)
			if err != nil{
				fmt.Println("match path error", err.Error())
				return false
			}
			return b
		},//自定义匹配
		Handler: func(w http.ResponseWriter, r *http.Request, next func(b bool), data *route.RouterData) {
			w.Write([]byte("user"))
			fmt.Println(r)
		},
	}
	_ = r
	a.AddRouter(r)
	r = &route.Route{
		Types:           []string{http.MethodGet},
		Path:            "/about",
		Match:           nil,
		Handler: func(w http.ResponseWriter, r *http.Request, next func(bool), data *route.RouterData) {
			w.Write([]byte("about"))
		},
	} //忽略
	a.AddRouter(r)
	r = &route.Route{
		Types: []string{http.MethodGet},
		Match: func(r *http.Request, data *route.RouterData) bool {
			b, err := regexp.MatchString(`^/user/([a-zA-Z0-9]+)/post/([a-zA-Z0-9]+)$`, r.URL.Path)
			if err != nil{
				fmt.Println("match path error", err.Error())
				return false
			}
			return b
		},
		Handler: func(w http.ResponseWriter, r *http.Request, next func(b bool), data *route.RouterData) {
			w.Write([]byte("user post"))
			fmt.Println(r)
		},
	} //忽略
	a.AddRouter(r)
	r = &route.Route{
		Types:           []string{http.MethodGet},
		Path:            "/old",// 直接匹配
		Match:           nil,
		Handler: func(w http.ResponseWriter, r *http.Request, next func(bool), data *route.RouterData) {
			w.Write([]byte("\nit's old"))
			data.Path = "/new"// 重定向
			next(false)
		},
	}
	a.AddRouter(r)
	r = &route.Route{
		Types:           []string{http.MethodGet},
		Path:            "/new",//重定向后继续匹配
		Match:           nil,
		Handler: func(w http.ResponseWriter, r *http.Request, next func(bool), data *route.RouterData) {
			w.Write([]byte("\nit's new"))
			next(false)
		},
	}
	a.AddRouter(r)
	r = &route.Route{
		Types:           []string{http.MethodGet},
		Path:            "/new",//重复匹配
		Match:           nil,
		Handler: func(w http.ResponseWriter, r *http.Request, next func(bool), data *route.RouterData) {
			w.Write([]byte("\nit's new2"))
			next(true)
		},
	}
	a.AddRouter(r)
	h := func(w http.ResponseWriter, r *http.Request, next func(bool), data *route.RouterData) {
		//w.Write([]byte("\nit's /*/b"))
		w.Write([]byte("matches is " + strings.Join(data.PathNodes, ",")))
		w.Write([]byte("\nid is " + data.MustGetValue("id").(string)))
		next(true)
	}
	a.Get("/*/b/:id/f", h)//胡乱匹配
	http.HandleFunc("/", a.Handler)
	http.ListenAndServe(":11041", nil)
}

