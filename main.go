package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"io"
	"io/ioutil"
	"os"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/k8scat/downhub/config"
	"github.com/k8scat/downhub/router"
	"github.com/k8scat/downhub/util"
)

var (
	cfgPath string
)

func main() {
	flag.StringVar(&cfgPath, "config", "./config.json", "config file path")
	flag.Parse()

	initLog()
	initConfig()
	initDingtalk()
	router.Run()
}

func initLog() {
	if err := os.MkdirAll("logs", 0744); err != nil {
		panic(err)
	}
	// write the logs to file and console at the same time
	logFile := fmt.Sprintf("logs/downhub_%s.log", time.Now().Format("20060102150405"))
	f, err := os.Create(logFile)
	if err != nil {
		if !os.IsExist(err) {
			panic(err)
		}
	}
	gin.DefaultWriter = io.MultiWriter(f, os.Stdout)
}

func initConfig() {
	b, err := ioutil.ReadFile(cfgPath)
	if err != nil {
		panic(err)
	}
	if err := json.Unmarshal(b, &config.Config); err != nil {
		panic(err)
	}
}

func initDingtalk() {
	if config.Config.DingtalkAccessToken == "" || config.Config.DingtalkSecret == "" {
		panic("dingtalk access_token or secret cannot be null")
	}
	var err error
	util.DingtalkClient, err = util.NewDingtalk(config.Config.DingtalkAccessToken, config.Config.DingtalkSecret)
	if err != nil {
		panic(err)
	}
}
