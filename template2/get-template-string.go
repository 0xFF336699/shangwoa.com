package template2

import (
	"bytes"
	ht "html/template"
	"sync"
	tt "text/template"
)

type templateCache struct {
	lock sync.RWMutex
	text map[string]*tt.Template
	html map[string]*ht.Template
}

var cache templateCache

func GetTextTemplate(str, name string, save bool) (*tt.Template, error) {
	if save {
		if t, ok := cache.text[name]; ok {
			return t, nil
		}
	}
	t, err := tt.New(name).Parse(str)
	if err != nil {
		return nil, err
	}
	if save {
		cache.text[name] = t
	}
	return t, nil
}

func GetTextTemplateBuffer(str, name string, data interface{}, save bool) (*bytes.Buffer, *tt.Template, error) {

	//t, err := tt.New(name).Parse(str)
	t, err := GetTextTemplate(str, name, save)
	if err != nil {
		return nil, nil, err
	}
	var b bytes.Buffer
	err = t.Execute(&b, data)
	return &b, t, err
}
