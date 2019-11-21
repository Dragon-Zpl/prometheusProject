package v1

import (
	"PrometheusProject/conf"
	"PrometheusProject/consul"
	"PrometheusProject/lib/stringi"
	"PrometheusProject/prometheus"
	"PrometheusProject/v1/form"
	"bytes"
	"encoding/json"
	"fmt"
	"github.com/astaxie/beego/logs"
	"github.com/lxzan/hasaki"
	"io/ioutil"
	"net/http"
	"os"
	"strings"
	"gopkg.in/gomail.v2"
	"sync"
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
	jsonValue := []byte(sendStr)
	resp, err := http.Post(url, "application/json", bytes.NewBuffer(jsonValue))
	if err != nil {
		fmt.Println(err)
	}
	fmt.Println(resp)
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

func RegisterOneVec(jobName, typ, labels, name string)  {
	if v, ok := prometheus.RegisterPromMap[jobName]; ok {
		var err error
		switch typ {
		case "counter":
			err = v.Vec.WithCounter(name, strings.Split(labels, ","))
		case "state":
			err = v.Vec.WithState(name, strings.Split(labels, ","))
		case "time":
			err = v.Vec.WithTimer(name, strings.Split(labels, ","))
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
			newProm.Vec.WithCounter(name, strings.Split(labels, ","))
		case "state":
			newProm.Vec.WithState(name, strings.Split(labels, ","))
		case "time":
			newProm.Vec.WithTimer(name, strings.Split(labels, ","))
		default:
			return
		}
		newProm.JobName = jobName
		newProm.Lock = sync.Mutex{}
		prometheus.RegisterPromMap[jobName] = newProm
	}
}