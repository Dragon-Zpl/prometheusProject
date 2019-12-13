package main

import (
	"PrometheusProject/conf"
	"PrometheusProject/prometheus"
	"PrometheusProject/route"
	v1 "PrometheusProject/v1"
)

func main() {
	conf.InitConf()
	go prometheus.GetPromHttp(":37983").ListenAndServe()
	// 注册已存在任务
	go v1.RegisterFromConsul()
	r := route.Router()
	r.Run(":37984")
}