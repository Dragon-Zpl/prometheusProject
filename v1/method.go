package v1

import (
	"PrometheusProject/conf"
	"PrometheusProject/consul"
	"PrometheusProject/lib/helper"
	"PrometheusProject/lib/stringi"
	"PrometheusProject/prometheus"
	"PrometheusProject/v1/form"
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/astaxie/beego/logs"
	"github.com/lxzan/hasaki"
	"io/ioutil"
	"net/http"
	"os"
	"strings"
	"gopkg.in/gomail.v2"
	"sync"
	"time"
)

func GetDirFileList(path string) map[string]struct{} {
	files, err := ioutil.ReadDir(path)
	res := make(map[string]struct{})
	if err != nil {
		return res
	}
	for _, file := range files {
		if file.IsDir() {
			continue
		}
		if strings.Contains(file.Name(), ".") {
			res[strings.Split(file.Name(), ".")[0]] = struct{}{}
		}
	}
	return res

}

func SignUpPrometheus()  {
	_, err := hasaki.Post(conf.PromeHotReloadCfg.PromeUrl, nil).Send(nil)

	if err != nil {
		logs.Error(err)
	}
	//fmt.Println(string(body.Body))

	_, err = hasaki.Post(conf.PromeHotReloadCfg.AlertUrl, nil).Send(nil)

	if err != nil {
		logs.Error(err)
	}
	//fmt.Println(string(body.Body))
	fmt.Printf("success")
}

var paramStr = `
{
    "msgtype": "{msgtype}",
    "markdown": {
        "title":"{title}",
        "text": "{text}",
    },
    "at": {
        "atMobiles": [
            {atMobiles}
        ],
        "isAtAll": {isAtAll}
    }
 }`





func SendDingDing(url string, data form.DingTalkRes) error {
	sendStr := stringi.Build(paramStr, stringi.Form{
		"msgtype": data.Msgtype,
		"title": data.Markdown.Title,
		"text": data.Markdown.Text,
		"atMobiles": strings.Join(data.At.AtMobiles, ","),
		"isAtAll": "false",
	})
	fmt.Println(sendStr)
	jsonValue := []byte(sendStr)
	resp, err := http.Post(url, "application/json", bytes.NewBuffer(jsonValue))
	if err != nil {
		fmt.Println(err)
	}
	respData, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		fmt.Println(err)
	}
	fmt.Println(helper.Date("Y-m-d H:i:s", time.Now().Unix()))
	fmt.Println(string(respData))
	return err
}

func PathExists(path string) bool{
	_, err := os.Stat(path)    //os.Stat获取文件信息
	if err != nil {
		if os.IsExist(err) {
			return true
		}
		return false
	}
	return true
}

func SendEmail(body, subject string, receive []string)  {
	m := gomail.NewMessage()

	m.SetHeader("From", conf.Email.Sender)
	m.SetHeader("To", receive...)
	m.SetHeader("Subject", subject)
	m.SetBody("text/html", body)
	d := gomail.NewDialer(conf.Email.Host, conf.Email.Port, conf.Email.Sender, conf.Email.Password)

	// Send the email to Bob, Cora and Dan.
	if err := d.DialAndSend(m); err != nil {
		logs.Error(err)
	}
	logs.Error("send email success")
}

func RegisterFromConsul()  {
	data := consul.GetKeyData("jobName")
	if len(data) == 0 {
		return
	}
	jobNames := make([]string, 0)
	err := json.Unmarshal(data, &jobNames)
	if err != nil {
		return
	}
	for _, jobName := range jobNames {
		jobDatas := consul.GetJobNameData(jobName)
		for _, jobData := range jobDatas {
			RegisterOneVec(jobName, jobData.Type, strings.Join(jobData.Labels, ","), jobData.Name)
			prometheus.VecMap[jobData.Name] = len(jobData.Labels)
		}
	}
}

func RegisterOneVec(jobName, typ, labels, name string, needReg ...bool)  {
	if v, ok := prometheus.RegisterPromMap[jobName]; ok {
		var err error
		switch typ {
		case "counter":
			err = v.Vec.WithCounter(name, strings.Split(labels, ","), needReg...)
		case "state":
			err = v.Vec.WithState(name, strings.Split(labels, ","), needReg...)
		case "time":
			err = v.Vec.WithTimer(name, strings.Split(labels, ","), needReg...)
		default:
			return
		}
		if err != nil {
			return
		}
	} else {
		newProm := prometheus.NewRegisterProm()
		newProm.Vec = prometheus.New()
		switch typ {
		case "counter":
			newProm.Vec.WithCounter(name, strings.Split(labels, ","), needReg...)
		case "state":
			newProm.Vec.WithState(name, strings.Split(labels, ","), needReg...)
		case "time":
			newProm.Vec.WithTimer(name, strings.Split(labels, ","), needReg...)
		default:
			return
		}
		newProm.JobName = jobName
		newProm.Lock = sync.Mutex{}
		prometheus.RegisterPromMap[jobName] = newProm
	}
}

func RegisterVecService(input form.RegisterVecForm) error {
	if input.JobName == "jobName" {
		return errors.New("jobName 存在请更换")
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
			return errors.New("不存在该类型指标")
		}
		if err != nil {
			return err
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
			return errors.New("不存在该类型指标")
		}
		if err != nil {
			return err
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
		if err := consul.SetKeyValue("jobName", jobNames); err != nil {
			return err
		}
	}
	consulData := consul.ConsulKVData{
		Name:   input.Name,
		Labels: strings.Split(input.Lables, ","),
		Type:   input.Typ,
	}

	if err := consul.SetJobNameValue(input.JobName, consulData); err != nil {
		return err
	}

	return nil
}

var TutuTaskData map[string]float64

func InitTutuConsul()  {
	TutuTaskData = make(map[string]float64)

	keys := consul.GetAllKey(conf.Consul.DirName + "/"  +conf.Consul.Prefix)
	var value float64
	for _, key := range keys {
		json.Unmarshal(consul.GetKeyData(key), &value)
		TutuTaskData[key] = value
		name := strings.Split(key, "/")[1]
		RegisterOneVec(conf.Consul.DirName, "state", "jobName", name, true)
		prometheus.PrometheusOpeartor(conf.Consul.DirName, name, value, []string{name}, prometheus.Set)
	}
}



func ListenConsul() {
	t := time.NewTimer(60 * time.Second)
	var value float64
	for   {
		select {
		case <-t.C:
			keys := consul.GetAllKey(conf.Consul.DirName + "/"  +conf.Consul.Prefix)
			for _, key := range keys {
				json.Unmarshal(consul.GetKeyData(key), &value)
				if v, ok := TutuTaskData[key]; ok {
					if v != value {
						name := strings.Split(key, "/")[1]
						prometheus.PrometheusOpeartor(conf.Consul.DirName, name, value, []string{name}, prometheus.Set)
					}
				} else {
					TutuTaskData[key] = value
					name := strings.Split(key, "/")[1]
					RegisterVecService(form.RegisterVecForm{
						JobName: conf.Consul.DirName,
						Name:    name,
						Lables:  "jobName",
						Typ:     "state",
					})
					prometheus.PrometheusOpeartor(conf.Consul.DirName, name, value, []string{name}, prometheus.Set)
				}
			}

			t.Reset(60 * time.Second)
		}
	}
}