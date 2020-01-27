package route

import (
	"net/http"
	"regexp"
	"strings"
)
const asterisk = "*"
const sharp = ":"
//const urlCharacters = `[-A-Za-z0-9+&@#%=~_|!:,.;]+`
const urlCharacters = `.*`
//const urlCharactersEnd = `[-A-Za-z0-9+&@#%=~_|!:,.;]+$`
const urlCharactersEnd = `.*$`
const urlCharactersAll = `[-A-Za-z0-9+&@#/%=~_|!:,.;]+`
//const urlCharactersAllEnd = `[-A-Za-z0-9+&@#/%=~_|!:,.;]*`
const urlCharactersAllEnd = `.*$`
type match func(r *http.Request, daata *RouterData)bool
type handler func(w http.ResponseWriter, r *http.Request, next func(bool),  data *RouterData)
// 如果某个步骤解析了部分数据可以放进来，设计思路为高内聚场景，路由之间基本已知互相的存在，这样可以节省计算
type matchNode struct{
	Type string // only * and :
	Index int
	Key string
}

type pathMatch struct{
	Nodes     []*matchNode
	RegString string
	Reg       *regexp.Regexp
}

type RouterData struct{
	R          *http.Request
	Path       string
	PathNodes []string // http://x.com/*/b/*/e 两个星号匹配的放在这个数组里
	PathParams map[string]string // http://x.com/:user/:post 这个对象里存放user,post
	Form       map[string]interface{}
	Body map[string]interface{}
	Queries    map[string]interface{}
	Extra      interface{} // 可传递next路由使用上一个路由产生的数据
}

func (d *RouterData) MustGetValue(key string)(value interface{})  {
	value, _ = d.GetValue(key)
	return
}

func (d *RouterData) GetValue(key string) (value interface{}, bl bool)  {
	if d.PathParams != nil {
		if value, bl = d.PathParams[key]; bl{
			return
		}
	}
	if d.Form != nil {
		if value, bl = d.Form[key]; bl{
			return
		}
	}
	if d.Queries != nil {
		if value, bl = d.Queries[key]; bl{
			return
		}
	}
	return
}

type Route struct{
	Types     []string // http.MethodGet
	Path      string
	Match     match
	pathMatch *pathMatch
	Handler   handler
}

type App struct{
	trees map[string][]*Route
}

func NewApp() (app *App) {
	return &App{ trees: map[string][]*Route{}}
}

func (app *App)Handler(w http.ResponseWriter, r *http.Request)  {
	runRouters(app, w, r)
}

func (app *App) AddRouter(router *Route) {
	if len(router.Types) == 0{
		panic("router types length cann't be zero")
	}
	for _, v := range router.Types{
		app.trees[v] = append(app.trees[v], router)
	}
}

func (app *App) Get(pattern string, handler handler) {
	app.AddRouter(createRouter(pattern, handler, []string{http.MethodGet}))
}
func createRouter(pattern string, handler handler, types []string) (router *Route) {
	ps := strings.Split(pattern, "/")
	pm := &pathMatch{
		Nodes:     []*matchNode{},
		RegString: "",
		Reg:       nil,
	}
	for i, p := range ps{
		if p == ""{
			continue
		}
		pm.RegString += "/"
		if strings.Index(p, asterisk) > -1{
			if i == len(ps) - 1{
				p = "(" + strings.Replace(p, asterisk, urlCharactersAllEnd, -1) + ")"
			}else{
				p = "(" + strings.Replace(p, asterisk, urlCharacters, -1) + ")"
			}
			pm.Nodes = append(pm.Nodes, &matchNode{
				Type:  asterisk,
				Index: i,
				Key:   "",
			})
		}else if strings.Index(p, sharp) > -1{
			key := strings.Split(p, sharp)[1]
			pm.Nodes = append(pm.Nodes, &matchNode{
				Type:  sharp,
				Index: i,
				Key:   key,
			})
			if i == len(ps) - 1{
				p = "(" + urlCharactersEnd + ")"
			}else{
				p = "(" + urlCharacters + ")"
			}
		}else{
			if i == len(ps) - 1{
				p += "$"
			}
		}
		pm.RegString += p
	}
	pm.Reg = regexp.MustCompile(pm.RegString)
	router = &Route{
		Types:     types,
		Path:      "",
		Match:     nil,
		pathMatch: nil,
		Handler:   handler,
	}
	if len(pm.Nodes) > 0{
		router.pathMatch = pm
	}else{
		router.Path = pattern
	}
	return
}
func matchPath(p string, pm *pathMatch, data *RouterData) (bl bool)  {
	bl, _ = regexp.MatchString(pm.RegString, p)
	ps := pm.Reg.FindStringSubmatch(p)
	if len(ps) - 1 != len(pm.Nodes) {
		return false
	}
	for i := 0; i < len(pm.Nodes); i++{
		mn := pm.Nodes[i]
		if mn.Type == sharp{
			if data.PathParams == nil{
				data.PathParams = map[string]string{}
			}
			data.PathParams[mn.Key] = ps[i + 1]
		}else{
			data.PathNodes = append(data.PathNodes, ps[i + 1])
		}
	}
	return
}

func runRouters(app *App, w http.ResponseWriter, r *http.Request) {
	handled := false
	index := 0
	data := &RouterData{
		R: r,
		Path:r.URL.Path,
	}
	routers := app.trees[r.Method]
	var next func(bool)
	next = func(end bool) {
		if end || index >= len(routers){
			return
		}
		router := routers[index]
		index ++
		if ((router.Path != "" && router.Path == data.Path) ||
				(router.Match != nil && router.Match(r, data)) ||
				(router.pathMatch != nil && matchPath(data.Path, router.pathMatch, data))) {
			handled = true
 			router.Handler(w, r, next, data)
		}else{
			next(false)
		}
		return
	}
	next(false)
	if !handled{
		w.WriteHeader(404)
	}
	return
}