package template2

import (
	ht "html/template"
	tt "text/template"
)

func init() {
	cache.text = map[string]*tt.Template{}
	cache.html = map[string]*ht.Template{}
}
