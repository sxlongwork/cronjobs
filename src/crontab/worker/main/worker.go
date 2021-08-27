package main

import (
	"cronjobs/src/crontab/worker"
	"cronjobs/src/crontab/worker/config"
	"flag"
	"fmt"
)

func main() {

	var (
		err        error
		configPath string
	)
	flag.StringVar(&configPath, "config", "worker.json", "the path of config file")
	flag.Parse()

	// 初始化配置
	if err = config.InitConfig(configPath); err != nil {
		fmt.Println("load config failed.", err)
	}

	// 启动日志协程
	if err = worker.InitLogMgr(); err != nil {
		fmt.Println(err)
	}

	// 初始化执行器
	if err = worker.InitJobExcutor(); err != nil {
		fmt.Println(err)
	}
	// 初始化调度器
	worker.InitSchduleMgr()

	// 初始化jobMgr, 启动监听
	if err = worker.InitJobMgr(); err != nil {
		fmt.Println(err)
	}

	select {}
}
