package helper

import (
	"PrometheusProject/lib/stringi"
	"bytes"
"fmt"
"github.com/gin-gonic/gin"
"net/http"
"regexp"
"strconv"
"strings"
"time"
)

func PreNDayTime(n int) int64 {
	return time.Now().AddDate(0, 0, -n).Unix()
}

func PreMinuteTime(duration time.Duration) int64 {
	return time.Now().Add(-duration).Unix()
}

func ArrayStrToInt(in []string) (outs []int) {
	for _, str := range in {
		if out, err := strconv.Atoi(str); err == nil {
			outs = append(outs, out)
		}
	}
	return
}

func SameElementCount(s []int) map[int]int {
	m := make(map[int]int)
	for _, elem := range s {
		if _, ok := m[elem]; !ok {
			m[elem] = 1
		} else {
			m[elem] += 1
		}
	}
	return m
}

// 驼峰转下划线
func Camel2Underline(s string) string {
	re, _ := regexp.Compile("[A-Z]{1}")
	s = re.ReplaceAllStringFunc(s, func(s string) string {
		m := []byte(s)
		return "_" + string(bytes.ToLower(m[0:1]))
	})
	return s
}

// 获取这周开始的时间
func GetWeekStart() time.Time {
	t := time.Now()
	t = time.Date(t.Year(), t.Month(), t.Day(), 0, 0, 0, 0, t.Location())
	wd := t.Weekday()
	if wd == time.Monday {
		return t
	}
	offset := int(time.Monday - wd)
	if offset > 0 {
		offset -= 7
	}
	return t.AddDate(0, 0, offset)
}

func Success() (int, interface{}) {
	return http.StatusOK, gin.H{
		"status": gin.H{
			"code":    0,
			"message": "success",
			"time": Date("Y-m-d H:i:s", time.Now().Unix()),
			"accessTokenState": "keep",
		},
		"data": gin.H{},
	}
}

func SuccessWithDate(data interface{}) (int, interface{}) {
	return http.StatusOK, gin.H{
		"status": gin.H{
			"code":    0,
			"message": "success",
			"time": Date("Y-m-d H:i:s", time.Now().Unix()),
			"accessTokenState": "keep",
		},
		"data":    data,
	}
}

func SuccessWithDataList(datalist interface{}) (int, interface{}) {
	return http.StatusOK, gin.H{
		"status": gin.H{
			"code":    0,
			"message": "success",
			"time": Date("Y-m-d H:i:s", time.Now().Unix()),
			"accessTokenState": "keep",
		},
		"data": gin.H{
			"dataList": datalist,
		},
	}
}

func SuccessWithDataAndPage(datalist interface{}, pageInfo interface{}) (int, interface{}) {
	return http.StatusOK, gin.H{
		"status": gin.H{
			"code":    0,
			"message": "success",
			"time": Date("Y-m-d H:i:s", time.Now().Unix()),
			"accessTokenState": "keep",
		},
		"data": map[string]interface{}{
			"dataList": datalist,
			"pageInfo": pageInfo,
		},
	}
}

func Date(format string, timestamp ...int64) string {
	var ts = time.Now().Unix()
	if len(timestamp) > 0 {
		ts = timestamp[0]
	}
	var t = time.Unix(ts, 0)
	Y := strconv.Itoa(t.Year())
	m := fmt.Sprintf("%02d", t.Month())
	d := fmt.Sprintf("%02d", t.Day())
	H := fmt.Sprintf("%02d", t.Hour())
	i := fmt.Sprintf("%02d", t.Minute())
	s := fmt.Sprintf("%02d", t.Second())

	format = strings.Replace(format, "Y", Y, -1)
	format = strings.Replace(format, "m", m, -1)
	format = strings.Replace(format, "d", d, -1)
	format = strings.Replace(format, "H", H, -1)
	format = strings.Replace(format, "i", i, -1)
	format = strings.Replace(format, "s", s, -1)
	return format
}

const ONEDAY = 60*60*24
func GetDay(time int64) string {
	if time == 0 {
		return "0"
	}
	return stringi.ToString(time/ONEDAY)
}


func GetPageInfo(page, pageCount int) (int, int) {
	if page == 0 {
		page = 1
	}
	if pageCount == 0 {
		pageCount = 5
	}

	return page, pageCount
}

func Fail(message string) (int, interface{}) {
	return http.StatusOK, gin.H{
		"status": gin.H{
			"code":    -1,
			"message": message,
			"time": Date("Y-m-d H:i:s", time.Now().Unix()),
			"accessTokenState": "keep",
		},
		"data": gin.H{},
	}
}