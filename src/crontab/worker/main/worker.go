package main

import (
	"cronjobs/src/crontab/worker"
	"cronjobs/src/crontab/worker/config"
	"flag"
	"log"
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
		log.Println(err)
	}

	// worker注册
	if err = worker.InitRegister(); err != nil {
		log.Println(err)
	}

	// 启动日志协程
	if err = worker.InitLogMgr(); err != nil {
		log.Println(err)
	}

	// 初始化任务执行器
	if err = worker.InitJobExcutor(); err != nil {
		log.Println(err)
	}
	// 初始化任务调度器
	if err = worker.InitSchduleMgr(); err != nil {
		log.Println(err)
	}

	// 初始化jobMgr, 启动监听
	if err = worker.InitJobMgr(); err != nil {
		log.Println(err)
	}

	select {}
}
