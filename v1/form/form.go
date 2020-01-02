package form



type AddRulesForm struct {
	Name string `form:"name" binding:"required"`
	AlertName string `form:"alertName" binding:"required"`
	Condition string `form:"condition" binding:"required"`
	Time string `form:"time" binding:"required"`
	Severity string `form:"severity" binding:"required"`
	Summary string `form:"summary" binding:"required"`
	Description string `form:"description" binding:"required"`
	DingUrl string `form:"dingUrl"`
	IsAtAll string `form:"isAtAll"`
	AtList string `form:"atList"`
	SendType string `form:"sendType"`
}


type DeleteRulesForm struct {
	Name string `form:"name" binding:"required"`
	AlertName string `form:"alertName"`
}

type MarkdownData struct {
	Title string `json:"title"`
	Text string `json:"text"`
}

type DingTalkAt struct {
	AtMobiles []string `json:"atMobiles"`
	IsAtAll bool `json:"isAtAll"`
}

type DingTalkRes struct {
	Msgtype string `json:"msgtype"`
	Markdown *MarkdownData `json:"markdown"`
	At *DingTalkAt `json:"at"`
}

type GetOneRulesForm struct {
	Name string `form:"name" binding:"required"`
	AlertName string `form:"alertName"`
}

type RegisterForm struct {
	ServerId string `form:"serverId"`
	ServerIp string `form:"serverIp"`
	ServerTag string `form:"serverTag"`
	ServerPath string `form:"serverPath"`
	ServerPort int `form:"serverPort"`
}

type DeRegisterForm struct {
	ServerIp string `form:"serverIp"`
	ServerId string `form:"serverId"`
}

type RegisterVecForm struct {
	JobName string `form:"jobName"`
	Name string `form:"name"`
	Lables string `form:"lables"`
	Typ    string `form:"type"`
}


type PromoperationForm struct {
	JobName string `form:"jobName"`
	Name string `form:"name"`
	Operation string `form:"operation"`
	Value float64 `form:"value"`
	Labels string `form:"labels"`
}

type UnRegisterVecForm struct {
	JobName string `form:"jobName"`
	Name string `form:"name"`
	Typ    string `form:"type"`
}

type CronWebHookForm struct {
	TaskId string `json:"task_id"`
	TaskName string `json:"task_name"`
	Status string `json:"status"`
	Result string `json:"result"`
	DingUrl string `json:"ding_url"`

}