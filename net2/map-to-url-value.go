package net2

import (
	"fmt"
	"net/url"
)

func MapToUV(m map[string]interface{}) url.Values {
	uv := url.Values{}
	for k, v := range m{
		uv.Set(k, fmt.Sprintf("%v", v))
	}
	return uv
}
