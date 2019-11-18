package main

import (
	"PrometheusProject/conf"
	"PrometheusProject/prometheus"
	"PrometheusProject/route"
	v1 "PrometheusProject/v1"
)

func main() {
	conf.InitConf()
	go prometheus.GetPromHttp(":6668").ListenAndServe()
	go v1.RegisterFromConsul()
	r := route.Router()
	r.Run(":6667")
}