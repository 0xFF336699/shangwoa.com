package route

import (
	"fmt"
	"io/ioutil"
	"mime"
	"net/http"
	"os"
	"path"
	"path/filepath"
	"regexp"
	"strconv"
	"strings"
	"sync"
	"time"
)
const asterisk = "*"
const sharp = ":"
const urlCharacters = `.*`
const urlCharactersEnd = `.*$`
const urlCharactersAll = `[-A-Za-z0-9+&@#/%=~_|!:,.;]+`
const urlCharactersAllEnd = `.*$`
type match func(r *http.Request, daata *RouterData)bool
// @next 建议每次经过路由，在末尾都调用next，并明确指示是否要结束，这样可以方便查看一个请求处理的耗时
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
	// 静态文件处理，runRouters方法里 如果请求是get方法，两个条件匹配后会在静态文件里查找
	//  handled==false，没有路由处理会进入静态文件查找
	// handled==true &&  IsStatic == true，有路由处理过，例如cookie和或者urlrewrite，但是没有输出，也会进行静态文件查找
	IsStatic bool
	Preffix  string
	NoCache bool
	IsFile bool
	ContentType string
}

func (d *RouterData) MustGetValue(key string)(value interface{})  {
	value, _ = d.GetValue(key)
	return
}

func (d *RouterData) MustGetStringValue(key string) string {
	v, b := d.GetValue(key)
	if !b{
		return ""
	}
	return v.(string)
}

func (d *RouterData) MustGetBoolValue(key string) bool {
	v, b := d.GetValue(key)
	if !b {
		return false
	}
	return v == "true" || v == "1"
}

func (d *RouterData) MustGetIntValue(key string) int {
	v, b := d.GetValue(key)
	if !b{
		return 0
	}
	n, e := strconv.Atoi(v.(string))
	if e != nil{
		return 0
	}
	return n
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
	values := d.R.URL.Query()
	v := values[key]
	if len(v) > 0{
		value = v[0]
		if d.Queries == nil{
			d.Queries = map[string]interface{}{}
		}
		d.Queries[key] = value
		bl = true
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
// 如果需要输出header的话，就需要把排序放在前面，注意 next(true)，这里的true是isEnd，是不再继续向下匹配的意思
// 多次匹配路由后，如果多次出现*星号匹配，则会累积，例如/* /b/* /*/b/* 在/*/b/*路由的时候data.PathNodes会有4个内容，所以 如果需要正确匹配就得用命名的方式:id :name这样
// @pattern /* 匹配 /任意后续路径和文件名，如 /a /a/ /a/b /a/b/c.html
// @pattern /a/* 匹配 /a/任意后续路径和文件名，如 /a/b /a/b/c.html
// 上面就是说 最后一个字符如果是星号，就匹配后面所有剩余部分，无论多少层路径和以及文件名
// @pattern /a/*/f 匹配 /a/任意字符/f
// @pattern /:name/*/:id 匹配 /任意字符被命名为name/任意字符/任意字符被命名为id
func (app *App) Get(pattern string, handler handler) {
	app.AddRouter(createRouter(pattern, handler, []string{http.MethodGet}))
}
func (app *App) OPTIONS(pattern string, handler handler) {
	app.AddRouter(createRouter(pattern, handler, []string{http.MethodOptions}))
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
	startTime := time.Now().Unix()
	handled := false
	index := 0
	data := &RouterData{
		R: r,
		Path:r.URL.Path,
	}
	complete := false
	routers := app.trees[r.Method]
	var next func(bool)
	next = func(end bool) {
		if end{
			complete = true
		}
		if end || index >= len(routers){
			fmt.Println("本次请求处理总耗时", r.RequestURI, data.Path, time.Now().Unix() - startTime)
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

	if r.Method == http.MethodGet && (!complete || (handled && data.IsStatic) || (!handled && data.IsStatic == false)){
		cf := GetFile(data.Preffix, data.Path, !data.NoCache)
		if cf != nil{
			handled = true
			if data.ContentType != ""{
				w.Header().Set("Content-Type", data.ContentType)
			}else if cf.ctype != ""{
				w.Header().Set("Content-Type", cf.ctype)
			}
			w.Header().Set("Content-Length", strconv.Itoa(len(cf.bs)))
			w.Header().Set("X-Upyun-Content-Length", strconv.Itoa(len(cf.bs)))
			writen := false
			if data.IsFile{
				//http.ServeFile(w, r, )
				fp := getFilePath(data.Preffix, data.Path)
				if fp != ""{
					http.ServeFile(w, r, fp)
					writen = true
				}
			}
			if !writen{
				w.Write(cf.bs)
			}
			fmt.Println("write content")
		}

	}
	if !handled{
		w.WriteHeader(404)
		return
	}
	return
}

type StaticDirectory struct {
	Directory string // 物理目录
	Path string // 映射路径前缀
}
type CachedFile struct{
	bs []byte
	ctype string
}
var staticDirectories []*StaticDirectory
var cachedFiles = &_StringICachedFileSyncMap{
}
var fileHandlers = map[string]func([]byte)[]byte{}

type _StringICachedFileSyncMap struct {
	Map sync.Map
}

func (this *_StringICachedFileSyncMap) Load(key string)(v *CachedFile, ok bool)  {
	i, ok := this.Map.Load(key)
	if !ok{
		return
	}
	v = i.(*CachedFile)
	return
}

func (this *_StringICachedFileSyncMap) LoadOrStore(key string, value *CachedFile)(v *CachedFile, ok bool)  {
	i, ok := this.Map.LoadOrStore(key, value)
	if !ok{
		return
	}
	v = i.(*CachedFile)
	return
}

func (this *_StringICachedFileSyncMap) Range(f func(key string, value *CachedFile) bool )  {
	this.Map.Range(func(key, value interface{}) bool {
		k := key.(string)
		v := value.(*CachedFile)
		return f(k, v)
	})
	return
}



func SetStaticDirectories(directories []*StaticDirectory)  {
	staticDirectories = directories
}

func AddFileHandler(p string, handler func([]byte) []byte) {
	fileHandlers[p] = handler
}

func getFilePath(p, f string)( s string)  {

	for i := 0; i < len(staticDirectories); i ++ {
		if staticDirectories[i].Path != p {
			continue
		}
		d := path.Join(staticDirectories[i].Directory, f)
		if fileExist(d){
			return  d
		}

	}
	return
}

func fileExist(path string) bool {
	_, err := os.Stat(path)    //os.Stat获取文件信息
	if err != nil {
		if os.IsExist(err) {
			return true
		}
		return false
	}
	return true
}
func GetFile(p, f string, cache bool) (cf *CachedFile) {
	var err error
	if cache{
		cf, hasFile := cachedFiles.Load(path.Join(p, f))
		if hasFile{
			return cf
		}
	}
	for i := 0; i < len(staticDirectories); i ++{
		if staticDirectories[i].Path != p{
			continue
		}
		d := path.Join(staticDirectories[i].Directory, f)
		d = strings.ReplaceAll(d, "/", string(os.PathSeparator))
		if _, err = os.Stat(d); os.IsNotExist(err) {
			fmt.Println("find file has error", p, f, err.Error())
			continue
		}
		bs, err := ioutil.ReadFile(d)
		if err != nil {
			fmt.Println("read file has error", p, f, err.Error())
			continue
		}
		if f, ok := fileHandlers[path.Join(p, f)]; ok{
			bs = f(bs)
		}

		ctype := mime.TypeByExtension(filepath.Ext(d))
		cf = &CachedFile{
			bs:    bs,
			ctype: ctype,
		}
		if cache{
			cachedFiles.Map.Store(path.Join(p, f), cf)
		}
		return
	}
	return
}