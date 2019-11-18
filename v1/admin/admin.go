package admin

import (
	"PrometheusProject/conf"
	"PrometheusProject/consul"
	"PrometheusProject/lib/helper"
	"PrometheusProject/lib/stringi"
	"PrometheusProject/prometheus"
	"PrometheusProject/v1"
	"PrometheusProject/v1/form"
	"encoding/json"
	"fmt"
	"github.com/astaxie/beego/logs"
	"github.com/gin-gonic/gin"
	"io/ioutil"
	"os"
	"strings"
	"sync"
)

func Mapping(prefix string, app *gin.Engine) {
	admin := app.Group(prefix)
	// 报警处理接口
	admin.POST("/webhook", PrometheusWebHook)
	// 报警规则添加
	admin.POST("/rules/add", PrometheusAddRules)
	// 报警规则删除
	admin.POST("/rules/delete", PrometheusAddDelete)
	// 向prometheus添加服务(target)
	admin.POST("/server/register", PrometheusRegister)
	// 注销
	admin.POST("/server/deregister", PrometheusDeRegister)
	// 注册数据指标
	admin.POST("/prom/register", RegisterVec)
	// 操作指标(增减)
	admin.POST("/prom/operation", Promoperation)
	// 删除数据指标
	admin.POST("/prom/unregister", UnRegisterVec)
}

func PrometheusWebHook(ctx *gin.Context)  {
	resData := make(map[string]interface{})
	data, _ := ioutil.ReadAll(ctx.Request.Body)
	err := json.Unmarshal([]byte(data), &resData)
	if err != nil {
		logs.Error(err)
	}
	var dingStr string
	kvModel := "{key}: {value} \n"
	notKvModel := "# {key}: \n ##### {value} \n"
	var output form.DingTalkRes
	output.Msgtype = "markdown"
	output.Markdown = new(form.MarkdownData)
	output.At = new(form.DingTalkAt)
	var url string
	var sendType string
	if alerts, ok := resData["alerts"].([]interface{}); ok{
		for k, v := range alerts[0].(map[string]interface{}) {
			if labels, ok := v.(map[string]interface{}); ok{
				kvStr := stringi.Build(kvModel, stringi.Form{
					"key": k,
					"value": "",
				})
				dingStr += "# " + kvStr
				for k, v := range labels {
					switch k {
					case "dingUrl":
						url = v.(string)
						continue
					case "isAtAll":
						if v == "true"{
							output.At.IsAtAll = true
						} else {
							output.At.IsAtAll = false
						}
						continue
					case "atList":
						if v != "" {
							output.At.AtMobiles = strings.Split(v.(string), ",")
							continue
						}
					case "alertname":
						output.Markdown.Title = v.(string)
					case "sendType":
						sendType = v.(string)
					}

					kvStr := stringi.Build(kvModel, stringi.Form{
						"key": k,
						"value": v.(string),
					})
					dingStr += "##### " +  kvStr
				}
			} else if value, ok := v.(string); ok{
				if k == "generatorURL" || k == "fingerprint" {
					continue
				}
				kvStr := stringi.Build(notKvModel, stringi.Form{
					"key": k,
					"value": value,
				})
				dingStr += kvStr
			}
		}
	}
	output.Markdown.Text = dingStr
	if url != "" && (sendType == "" || sendType == "ding"){
		v1.SendDingDing(url, output)
	} else if sendType == "email" {
		v1.SendEmail(strings.Replace(strings.Replace(output.Markdown.Text, "#", "", -1), "\n", "<br>", -1), output.Markdown.Title, output.At.AtMobiles)
	} else {
		ctx.JSON(helper.Fail("不存在该发送方式"))
	}
	ctx.JSON(helper.Success())
	return
}

func PrometheusAddRules(ctx *gin.Context)  {
	var input form.AddRulesForm
	if err := ctx.ShouldBind(&input); err != nil {
		ctx.JSON(helper.Fail("参数校验失败"))
		return
	}
	fileList := v1.GetDirFileList(conf.PromethuesPath.Path)

	if strings.Contains(input.Name, ".") || (!strings.Contains(input.Time, "s") && !strings.Contains(input.Time, "m") && !strings.Contains(input.Time, "d")) {
		ctx.JSON(helper.Fail("name can,t contain '.'!"))
		return
	}
	yamlFile, err := ioutil.ReadFile(conf.PromethuesPath.Path + "rules.yml")
	if err != nil {
		ctx.JSON(helper.Fail("file no exist"))
		return
	}
	if _, ok := fileList[input.Name]; !ok {
		data := stringi.Build(string(yamlFile), stringi.Form{
			"name": input.Name,
			"alertName": input.AlertName,
			"condition": input.Condition,
			"time": input.Time,
			"severity": input.Severity,
			"summary": input.Summary,
			"description": input.Description,
			"dingUrl": input.DingUrl,
			"isAtAll": input.IsAtAll,
			"alList": input.AtList,
			"sendType": input.SendType,
		})
		modelFilePath := conf.PromethuesPath.Path + input.Name + ".rules"
		fmt.Printf(modelFilePath)
		err = ioutil.WriteFile(modelFilePath, []byte(data), 0644)
		if err != nil {
			ctx.JSON(helper.Fail(err.Error()))
			return
		}

	} else {
		oldFile, err := ioutil.ReadFile(conf.PromethuesPath.Path + input.Name +".rules")

		if err != nil {
			ctx.JSON(helper.Fail("file no exist"))
			return
		}

		var isExist bool
		var alertIndex int
		yamlSlice := strings.Split(string(oldFile), "\n")
		for index, line := range yamlSlice {
			if strings.Contains(line, input.AlertName) {
				isExist = true
				alertIndex = index
				break
			} else {
				isExist = false
			}
		}
		data := stringi.Build(strings.Join(strings.Split(string(yamlFile), "\n")[3:], "\n"), stringi.Form{
			"alertName": input.AlertName,
			"condition": input.Condition,
			"time": input.Time,
			"severity": input.Severity,
			"summary": input.Summary,
			"description": input.Description,
			"dingUrl": input.DingUrl,
			"isAtAll": input.IsAtAll,
			"alList": input.AtList,
			"sendType": input.SendType,
		})
		fileData := make([]string, 0)
		if !isExist{

			fileData = append(fileData, string(oldFile))

			fileData = append(fileData, data)


		} else {
			newFileData := make([]string, 0)
			for index, line := range yamlSlice {
				if index < alertIndex || index > alertIndex + 12 {
					newFileData = append(newFileData, line)
				}
			}

			fileData = append(fileData, strings.Join(newFileData, "\n"))
			fileData = append(fileData, data)

		}
		fileStr := strings.Join(fileData, "\n")
		err = ioutil.WriteFile(conf.PromethuesPath.Path + input.Name + ".rules", []byte(fileStr), 0644)
		if err != nil {
			ctx.JSON(helper.Fail(err.Error()))
		}
	}
	v1.SignUpPrometheus()
	ctx.JSON(helper.Success())
	return
}


func PrometheusAddDelete(ctx *gin.Context)  {
	var input form.DeleteRulesForm
	if err := ctx.ShouldBind(&input); err != nil {
		ctx.JSON(helper.Fail("参数校验失败"))
		return
	}
	fileName := input.Name + ".rules"

	if !v1.PathExists(conf.PromethuesPath.Path + fileName){
		ctx.JSON(helper.Fail("文件不存在"))
	}

	if input.AlertName == "" {
		err := os.Remove(conf.PromethuesPath.Path + fileName)
		if err != nil {
			ctx.JSON(helper.Fail("删除失败"))
			return
		}
	} else {
		yamlFile, err := ioutil.ReadFile(conf.PromethuesPath.Path + fileName)
		if err != nil {
			ctx.JSON(helper.Fail("删除失败"))
			return
		}
		yamlSlice := strings.Split(string(yamlFile), "\n")
		newFileData := make([]string, 0)
		var isExist bool
		var alertIndex int
		for index, line := range yamlSlice {
			if strings.Contains(line, input.AlertName) {
				isExist = true
				alertIndex = index
				break
			} else {
				isExist = false
			}
		}

		if isExist {
			for index, line := range yamlSlice {
				if index < alertIndex || index > alertIndex + 12 {
					newFileData = append(newFileData, line)
				}
			}
			err = os.Remove(conf.PromethuesPath.Path + fileName)
			fileStr := strings.Join(newFileData, "\n")
			err = ioutil.WriteFile(conf.PromethuesPath.Path + fileName, []byte(fileStr), 0644)
		} else {
			ctx.JSON(helper.Fail("alterName 不存在"))
			return
		}

	}
	v1.SignUpPrometheus()
	ctx.JSON(helper.Success())
	return

}

func PrometheusRegister(ctx *gin.Context)  {
	var input form.RegisterForm
	if err := ctx.ShouldBind(&input); err != nil {
		ctx.JSON(helper.Fail("参数校验失败"))
		return
	}

	if consul.IsExistService(input.ServerId, input.ServerIp) {
		ctx.JSON(helper.Fail("serverId 已存在"))
		return
	}

	err := consul.ConsulRegisterServer(input.ServerId, input.ServerIp, input.ServerIp, input.ServerTag, input.ServerPath, input.ServerPort)

	if err != nil {
		ctx.JSON(helper.Fail(err.Error()))
	}

	ctx.JSON(helper.Success())
	return
}

func PrometheusDeRegister(ctx *gin.Context)  {
	var input form.DeRegisterForm
	if err := ctx.ShouldBind(&input); err != nil {
		ctx.JSON(helper.Fail("参数校验失败"))
		return
	}
	if !consul.IsExistService(input.ServerId, input.ServerIp) {
		ctx.JSON(helper.Fail("serverId 不存在"))
		return
	}
	err := consul.ConsulDeRegisterServer(input.ServerId, input.ServerIp)
	if err != nil {
		ctx.JSON(helper.Fail(err.Error()))
	}

	ctx.JSON(helper.Success())
	return
}

func RegisterVec(ctx *gin.Context)  {
	var input form.RegisterVecForm
	if err := ctx.ShouldBind(&input); err != nil {
		ctx.JSON(helper.Fail("参数校验失败"))
		return
	}
	if input.JobName == "jobName" {
		ctx.JSON(helper.Fail("jobName 存在请更换"))
		return
	}
	if v, ok := prometheus.RegisterPromMap[input.JobName]; ok {
		var err error
		switch input.Typ {
		case "counter":
			err = v.Vec.WithCounter(input.Name, strings.Split(input.Lables, ","))
		case "state":
			err = v.Vec.WithState(input.Name, strings.Split(input.Lables, ","))
		case "time":
			err = v.Vec.WithTimer(input.Name, strings.Split(input.Lables, ","))
		default:
			ctx.JSON(helper.Fail("不存在该类型指标"))
			return
		}
		if err != nil {
			ctx.JSON(helper.Fail(err.Error()))
			return
		}
	} else {
		newProm := prometheus.NewRegisterProm()
		newProm.Vec = prometheus.New()
		var err error
		switch input.Typ {
		case "counter":
			err = newProm.Vec.WithCounter(input.Name, strings.Split(input.Lables, ","))
		case "state":
			err = newProm.Vec.WithState(input.Name, strings.Split(input.Lables, ","))
		case "time":
			err = newProm.Vec.WithTimer(input.Name, strings.Split(input.Lables, ","))
		default:
			ctx.JSON(helper.Fail("不存在该类型指标"))
			return
		}
		if err != nil {
			ctx.JSON(helper.Fail(err.Error()))
			return
		}
		newProm.JobName = input.JobName
		newProm.Lock = sync.Mutex{}
		prometheus.RegisterPromMap[input.JobName] = newProm
		jobNames := make([]string, 0)
		if data := consul.GetKeyData("jobName"); len(data) > 0{
			oldJobName := make([]string, 0)
			json.Unmarshal(data, &oldJobName)
			flag := 0
			for _, v := range oldJobName {
				if v == input.JobName {
					flag = 1
				}
			}
			if flag == 0 {
				jobNames = append(oldJobName, input.JobName)
			} else {
				jobNames = oldJobName
			}
		} else {
			jobNames = append(jobNames, input.JobName)
		}
		consul.SetKeyValue("jobName", jobNames)
	}
	consulData := consul.ConsulKVData{
		Name:   input.Name,
		Labels: strings.Split(input.Lables, ","),
		Type:   input.Typ,
	}

	consul.SetJobNameValue(input.JobName, consulData)
	ctx.JSON(helper.Success())

	return
}

func Promoperation(ctx *gin.Context)  {
	var input form.PromoperationForm
	if err := ctx.ShouldBind(&input); err != nil {
		ctx.JSON(helper.Fail("参数校验失败"))
		return
	}

	var opt prometheus.QueryType

	switch input.Operation {
	case "add":
		opt = prometheus.Add
	case "inc":
		opt = prometheus.Inc
	case "dec":
		opt = prometheus.Dec
	case "time":
		opt = prometheus.Timing
	case "set":
		opt = prometheus.Set
	}

	if _, ok := prometheus.RegisterPromMap[input.JobName]; !ok {
		ctx.JSON(helper.Fail("不存在该JobName,请先注册"))
		return
	}

	err := prometheus.PrometheusOpeartor(input.JobName, input.Name, input.Value, strings.Split(input.Labels, ","), opt)
	if err != nil {
		ctx.JSON(helper.Fail(err.Error()))
		return
	}

	ctx.JSON(helper.Success())
	return
}

func UnRegisterVec(ctx *gin.Context)  {
	var input form.UnRegisterVecForm

	if err := ctx.ShouldBind(&input); err != nil {
		ctx.JSON(helper.Fail("参数校验失败"))
		return
	}

	if v, ok := prometheus.RegisterPromMap[input.JobName]; ok {
		err := v.Vec.UnRegister(input.Typ, input.Name)
		if err != nil {
			ctx.JSON(helper.Fail(err.Error()))
			return
		}
	} else {
		ctx.JSON(helper.Fail("不存在该JobName"))
	}

	consul.DeleteJobNameData(input.JobName, input.Name)
	prometheus.DeleteVec(input.Name)
	ctx.JSON(helper.Success())
	return
}