package main

import (
	"PrometheusProject/conf"
	"PrometheusProject/consul"
	"PrometheusProject/prometheus"
	"sync"
)



func main() {
	conf.InitConf()
	n := prometheus.New()
	n.WithCounter("test", []string{"abc"})
	consul.SetKeyValue("test1", prometheus.RegisterProm{
		Vec:     n,
		JobName: "test2",
		Lock:    sync.Mutex{},
	})

}
