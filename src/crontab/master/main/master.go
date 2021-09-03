package main

import (
	"cronjobs/src/crontab/master"
	"cronjobs/src/crontab/master/config"
	"flag"
	"fmt"
)

func main() {
	var (
		err        error
		configPath string
	)
	flag.StringVar(&configPath, "config", "master.json", "the path of config file")
	flag.Parse()

	// 初始化配置
	if err = config.InitConfig(configPath); err != nil {
		fmt.Println("load config failed.", err)
	}

	// 初始化workerMgr
	if err = master.InitWorkerMgr(); err != nil {
		fmt.Println(err)
	}
	// 初始化jobMgr, etcd信息
	if err = master.InitJobMgr(); err != nil {
		fmt.Println(err)
	}
	// 初始化日志管理器
	if err = master.InitLogMgr(); err != nil {
		fmt.Println(err)
	}

	// 初始化apiserver
	if err = master.InitApiServer(); err != nil {
		fmt.Println("start api server failed", err)
	}
	select {}

}
