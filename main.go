package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"os"
	"sync"
	"time"

	"gorm.io/driver/mysql"

	_ "github.com/go-sql-driver/mysql"

	"gorm.io/gorm"
)

var (
	db        *gorm.DB
	DBNAME    string
	DBUSER    string
	DBPWD     string
	DBHOST    string
	SERVPPORT string
	DEBUGFLAG string
)

func init() {
	var err error
	//如 testdb
	DBNAME = os.Getenv("logwebhook_dbname")
	//如 localhost:6666
	DBHOST = os.Getenv("logwebhook_dbhost")
	DBUSER = os.Getenv("logwebhook_dbuser")
	DBPWD = os.Getenv("logwebhook_dbpwd")
	SERVPPORT = os.Getenv("logwebhook_servport")
	DEBUGFLAG = os.Getenv("logwebhook_debugflag")
	if DEBUGFLAG == "" {
		DEBUGFLAG = "debug"
	}

	if DBNAME == "" || DBHOST == "" || SERVPPORT == "" {
		panic("ENV NOT SETTED")
	}
	db, err = gorm.Open(mysql.New(mysql.Config{

		DSN: fmt.Sprintf("%s:%s@tcp(%s)/%s?charset=utf8&parseTime=True&loc=Local",
			DBUSER, DBPWD, DBHOST, DBNAME),
	}), &gorm.Config{})
	if err != nil {
		panic("panic: mysql connnection panic")
	}
	alertMap.existedMap = make(map[string]bool)
	alertMap.serialMap = make(map[string]int)
	alerts := make([]AlertLog, 0)
	db.Where("status = ? ", "firing").Find(&alerts)
	for _, a := range alerts {
		alertMap.existedMap[a.Alertname] = true
		alertMap.serialMap[a.Alertname] = a.Id
	}
}
func printMap() {
	for {
		fmt.Println("========existedMap========")
		fmt.Println(alertMap.existedMap)
		fmt.Println("========serialMap========")
		fmt.Println(alertMap.serialMap)
		time.Sleep(10 * time.Second)
	}
}
func main() {
	//go printMap()
	http.HandleFunc("/webhook", webhookHandle)
	http.ListenAndServe(SERVPPORT, nil)
}

type AlertLog struct {
	Id         int       `json:"id" gorm:"id"`
	Alertname  string    `json:"alertname" gorm:"alertname"`
	Count      int       `json:"count" gorm:"count"`
	Status     string    `json:"status" gorm:"status"`
	UpdateTime time.Time `json:"update_time" gorm:"update_time"`
	CreateTime time.Time `json:"create_time" gorm:"create_time"`
}

func (AlertLog) TableName() string {
	return "alertlog"
}

type AlertLogMap struct {
	sync.Mutex
	//告警存在标记
	existedMap map[string]bool
	//告警所在序列号(id)
	serialMap map[string]int
}

func (alertLog *AlertLogMap) IfExists(alertname string) (bool, int) {
	alertLog.Lock()
	defer alertLog.Unlock()
	return alertLog.existedMap[alertname], alertLog.serialMap[alertname]
}

func (alertLog *AlertLogMap) AddAlertflag(alertname string) {
	alertLog.Lock()
	defer alertLog.Unlock()
	alertLog.existedMap[alertname] = true
}
func (alertLog *AlertLogMap) AddAlertserial(alertname string, id int) {
	alertLog.Lock()
	defer alertLog.Unlock()
	alertLog.serialMap[alertname] = id
}

func (alertLog *AlertLogMap) DeleteAlert(alertname string) {
	alertLog.Lock()
	defer alertLog.Unlock()
	delete(alertLog.existedMap, alertname)
	delete(alertLog.serialMap, alertname)
}

var alertMap AlertLogMap

type alertInfo struct {
	CommonLabels struct {
		Alertname string `json:"alertname"`
	} `json:"commonLabels"`
	//resolved firing
	Status string `json:"status"`
	Alerts []struct {
		//resolved firing
		Status string `json:"status"`
		Labels struct {
			Alertname string `json:"alertname"`
			Name      string `json:"name"`
		} `json:"labels"`
	} `json:"alerts"`
}

func webhookHandle(w http.ResponseWriter, req *http.Request) {

	bs, err := ioutil.ReadAll(req.Body)
	if err != nil {
		fmt.Println("err:", err)
		return
	}
	if DEBUGFLAG == "debug" {
		fmt.Println(string(bs))
	}
	var alert alertInfo

	if err := json.Unmarshal(bs, &alert); err != nil {
		fmt.Println("err json:", err)
		return
	}
	switch alert.Status {
	case "resolved":
		handleResolved(alert.CommonLabels.Alertname)
	case "firing":
		handleFiring(alert.CommonLabels.Alertname)
	}

}
func handleResolved(alertname string) {
	exist, serialNum := alertMap.IfExists(alertname)
	if exist {
		db.Model(&AlertLog{}).Where("id = ?", serialNum).Updates(map[string]interface{}{
			"update_time": time.Now(),
			"status":      "resolved",
		})
	}
	alertMap.DeleteAlert(alertname)
}
func handleFiring(alertname string) {
	exist, serialNum := alertMap.IfExists(alertname)
	var alert AlertLog
	//如果不存在，那么新增
	if !exist {
		alertMap.AddAlertflag(alertname)
		alert = AlertLog{
			Alertname:  alertname,
			Count:      1,
			Status:     "firing",
			UpdateTime: time.Now(),
			CreateTime: time.Now(),
		}
		db.Create(&alert)
		alertMap.AddAlertserial(alertname, alert.Id)
	} else {
		db.First(&alert, serialNum)
		db.Model(&AlertLog{}).Where("id = ?", serialNum).Updates(map[string]interface{}{
			"count":       alert.Count + 1,
			"update_time": time.Now(),
		})
	}
}
