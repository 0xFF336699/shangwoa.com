package  monitor

import (
	"net/http"
	_ "net/http/pprof"
)

func StartMonitor(port string)  {
	http.ListenAndServe("0.0.0.0:" + port, nil)
}