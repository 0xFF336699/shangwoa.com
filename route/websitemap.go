package route

import (
	"mime"
	"net/http"
	"path"
	"path/filepath"
	"shangwoa.com/os2"
	"strings"
	"fmt"
)

func init() {
	mime.AddExtensionType(".apk", "application/vnd.android.package-archive")
}
type FilePath struct {
	//静态文件目录 例如 "F:\\work\\express\\web\\express-web-docs"
	Directory string `json:"directory"`
}
type Site struct{
	// 路径前缀，例如http://x.com/a/index.html
	// prefix就是/a
	Prefix    string `json:"prefix"`
	FilePaths []*FilePath
	// 需要重定向的文件，例如index.html，或者其它一些看似私有路径的url可以在这里映射到其它文件上去
	// 假设url路径为http://x.com/a/b/c.html，prefix为/a
	// 若需要把c.html映射到index.html，则map映射关系为{"/b/c.html":"/index.html"}
	RedirectMap map[string]string
	Hosts []string
}

type WebsiteConf struct {
	Map map[string]*Site
}
var websiteConf *WebsiteConf
type WebsiteConf2 struct {
	Map map[string] []*Site
}
var websiteConf2 *WebsiteConf2
var confPath string
func UpdateWebsiteConf() (err error) {
	websiteConf = new(WebsiteConf)
	err = os2.LoadFileToStruct(confPath, websiteConf)
	initWebsites()
	return err
}

func initWebsites() {
	websiteConf2 = &WebsiteConf2{
		Map: map[string][]*Site{},
	}
	for k, v := range websiteConf.Map{
		websiteConf2.Map[k] = []*Site{v}
		for _, host := range v.Hosts{
			if websiteConf2.Map[host] == nil{
				websiteConf2.Map[host] = []*Site{v}
			}else{
				websiteConf2.Map[host] = append(websiteConf2.Map[host], v)
			}
		}
	}
}
func SetWebsiteConf(p string) (err error) {
	confPath = p
	return UpdateWebsiteConf()
}

func ServeFileIfExist2(w http.ResponseWriter, r *http.Request, p string) (err error, handled bool) {
	err, realPath, hasFile := GetFilePathByUrlPath(p)
	if err != nil{
		return
	}
	if !hasFile {
		return
	}
	http.ServeFile(w, r, realPath)
	return nil, true
}

func ServeFileIfExist(w http.ResponseWriter, r *http.Request, p string) (err error, handled bool) {
	err, realPath, hasFile, hasSite := GetFilePathByHostsUrl(r.Host, p)
	if err != nil{
		return
	}
	if !hasFile && !hasSite{
		err, realPath, hasFile = GetFilePathByUrlPath(p)
	}
	if err != nil{
		return
	}
	if !hasFile {
		return
	}
	http.ServeFile(w, r, realPath)
	return nil, true
}

func GetFilePathByHostsUrl(host, p string) (err error, realPath string, hasFile, hasSite bool) {
	if websiteConf == nil{
		return
	}
	sites := websiteConf2.Map[host]
	if sites == nil || len(sites) == 0{
		hp := strings.Split(host, ".")
		hp[0] = "*"
		host = strings.Join(hp[:], ".")
		//return
	}
	sites = websiteConf2.Map[host]
	if sites == nil || len(sites) == 0{
		return
	}
	hasSite = true
	for _, site := range sites{
		err, realPath, hasFile = GetFileBySite(site, p)
		if hasFile{
			return
		}
	}
	return
}

func GetFileBySite(site *Site, fp string)  (err error, realPath string, hasFile bool) {
	fmt.Println("site is", site.Hosts[0], site.FilePaths[0].Directory, len(site.RedirectMap))
	if _, ok := site.RedirectMap[fp]; ok{
		fp = site.RedirectMap[fp]
	}
	for _, pathConf := range site.FilePaths{
		d := path.Join(pathConf.Directory, fp)
		d = filepath.FromSlash(d)
		fmt.Println("d is ", d)
		if fileExist(d){
			return  nil, d, true
		}
	}
	return
}
func GetFilePathByUrlPath(p string) (err error, realPath string, hasFile bool) {
	if websiteConf == nil{
		return
	}
	bs := []byte(p)
	if len(bs) < 3{
		// "/a/带前缀至少是3个字符，低于三个字符应该跳出
		return
	}
	bs = bs[1:]
	index := strings.Index(string(bs), "/")
	if index == -1{
		return
	}
	prefix := string(bs[:index])
	site, ok := websiteConf.Map[prefix]
	if !ok {
		return
	}
	fp := string(bs[index:])
	if _, ok = site.RedirectMap[fp]; ok{
		fp = site.RedirectMap[fp]
	}
	for _, pathConf := range site.FilePaths{
		d := path.Join(pathConf.Directory, fp)
		d = filepath.FromSlash(d)
		if fileExist(d){
			return  nil, d, true
		}
	}
	return
}
