package view

import (
	"PrometheusProject/conf"
	"PrometheusProject/consul"
	"PrometheusProject/lib/helper"
	"PrometheusProject/v1"
	"PrometheusProject/v1/form"
	"github.com/gin-gonic/gin"
	"io/ioutil"
	"strings"
)

func Mapping(prefix string, app *gin.Engine) {
	view := app.Group(prefix)
	view.GET("/rulesList", GetRulesList)
	view.GET("/rules", GetOneRules)
	view.GET("/consul/list", GetConsulServer)
}

//查看有哪些报警规则组
func GetRulesList(ctx *gin.Context)  {
	datas := v1.GetDirFileList(conf.PromethuesPath.Path)
	res := make([]string, 0)
	for data, _ := range datas {
		res =append(res, data)
	}
	ctx.JSON(helper.SuccessWithDate(datas))
	return
}

// 获取报警规则
func GetOneRules(ctx *gin.Context)  {
	var input form.GetOneRulesForm
	if err := ctx.ShouldBind(&input); err != nil {
		ctx.JSON(helper.Fail("参数校验失败"))
		return
	}

	fileName := input.Name + ".rules"
	yamlFile, err := ioutil.ReadFile(conf.PromethuesPath.Path + fileName)
	if err != nil {
		ctx.JSON(helper.Fail("file no exist"))
		return
	}
	if input.AlertName == "" {
		ctx.JSON(helper.SuccessWithDataList(string(yamlFile)))
		return
	} else {
		var isExist bool
		var alertIndex int
		input.AlertName = input.AlertName + "警告"
		yamlSlice := strings.Split(string(yamlFile), "\n")
		for index, line := range yamlSlice {
			if strings.Contains(line, input.AlertName) {
				isExist = true
				alertIndex = index
				break
			} else {
				isExist = false
			}
		}
		if !isExist {
			ctx.JSON(helper.Fail("alterName 不存在"))
			return
		}
		ctx.JSON(helper.SuccessWithDataList(strings.Join(yamlSlice[alertIndex: alertIndex + 12], "\n")))
		return
	}



}

// 查看注册的服务
func GetConsulServer(ctx *gin.Context)  {
	res := consul.GetConsulRegisterService()

	ctx.JSON(helper.SuccessWithDataList(res))
	return
}














